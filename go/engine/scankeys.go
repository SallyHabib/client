package engine

import (
	"github.com/keybase/client/go/libkb"
	"golang.org/x/crypto/openpgp"
)

// ScanKeys finds pgp decryption keys in SKB and also if there is
// one stored on the server.  It satisfies the openpgp.KeyRing
// interface.
//
// It also will find public pgp keys for signature verification.
//
// It is not an engine, but uses an engine and is used by engines,
// so has to be in the engine package.  It is a UIConsumer.
type ScanKeys struct {
	keys  openpgp.EntityList
	secui libkb.SecretUI
	idui  libkb.IdentifyUI
	opts  *TrackOptions
}

// enforce ScanKeys implements openpgp.KeyRing:
var _ openpgp.KeyRing = &ScanKeys{}

// NewScanKeys creates a ScanKeys type.  The user parameter is
// optional.
func NewScanKeys(u *libkb.User, secui libkb.SecretUI, idui libkb.IdentifyUI, opts *TrackOptions) (*ScanKeys, error) {
	sk := &ScanKeys{secui: secui, idui: idui, opts: opts}
	if u == nil {
		return sk, nil
	}

	// if user provided, then load their local keys, and their synced secret key:
	ring, err := G.LoadSKBKeyring(u.GetName())
	if err != nil {
		return nil, err
	}
	synced, err := u.GetSyncedSecretKey()
	if err != nil {
		return nil, err
	}
	if err := sk.extractKeys(ring, synced, secui); err != nil {
		return nil, err
	}
	return sk, nil
}

func (s *ScanKeys) Name() string {
	return "ScanKeys"
}

func (s *ScanKeys) RequiredUIs() []libkb.UIKind {
	return []libkb.UIKind{libkb.SecretUIKind}
}

func (s *ScanKeys) SubConsumers() []libkb.UIConsumer {
	return []libkb.UIConsumer{
		&PGPKeyfinder{},
	}
}

// Count returns the number of local keys available.
func (s *ScanKeys) Count() int {
	return len(s.keys)
}

// KeysById returns the set of keys that have the given key id.
func (s *ScanKeys) KeysById(id uint64) []openpgp.Key {
	memres := s.keys.KeysById(id)
	G.Log.Debug("ScanKeys:KeysById(%d) => %d keys match in memory", id, len(memres))
	if len(memres) > 0 {
		return memres
	}

	// no match, so use api server to find keys for this id
	list, err := s.scan(id)
	if err != nil {
		G.Log.Warning("error finding keys: %s", err)
		return nil
	}
	G.Log.Debug("ScanKeys:KeysById(%d) => %d keys found via api scan", id, len(list))
	return list.KeysById(id)
}

// KeysByIdAndUsage returns the set of keys with the given id
// that also meet the key usage given by requiredUsage.
// The requiredUsage is expressed as the bitwise-OR of
// packet.KeyFlag* values.
func (s *ScanKeys) KeysByIdUsage(id uint64, requiredUsage byte) []openpgp.Key {
	// first, check the keys we already extracted.
	memres := s.keys.KeysByIdUsage(id, requiredUsage)
	G.Log.Debug("ScanKeys:KeysByIdUsage(%d, %x) => %d keys match in memory", id, requiredUsage, len(memres))
	if len(memres) > 0 {
		return memres
	}

	// no match, so now lookup the user on the api server by the key id.
	list, err := s.scan(id)
	if err != nil {
		G.Log.Warning("error finding keys: %s", err)
		return nil
	}
	// use the list to find the keys correctly
	G.Log.Debug("ScanKeys:KeysByIdUsage(%d, %x) => %d keys found via api scan", id, requiredUsage, len(list))
	return list.KeysByIdUsage(id, requiredUsage)
}

// DecryptionKeys returns all private keys that are valid for
// decryption.
func (s *ScanKeys) DecryptionKeys() []openpgp.Key {
	G.Log.Debug("ScanKeys:DecryptionKeys() => %d keys available", len(s.keys))
	return s.keys.DecryptionKeys()
}

// extractKeys gets all the private pgp keys out of the ring and
// the synced key.
func (s *ScanKeys) extractKeys(ring *libkb.SKBKeyringFile, synced *libkb.SKB, ui libkb.SecretUI) error {
	if err := s.extractKey(synced, ui); err != nil {
		return err
	}

	for _, b := range ring.Blocks {
		if !libkb.IsPgpAlgo(b.Type) {
			continue
		}
		if err := s.extractKey(b, ui); err != nil {
			return err
		}
	}

	return nil
}

// extractKey gets the private key out of skb.  If it's a pgp key,
// it adds it to the keys stored in s.
func (s *ScanKeys) extractKey(skb *libkb.SKB, ui libkb.SecretUI) error {
	if skb == nil {
		return nil
	}
	k, err := skb.PromptAndUnlock("pgp decrypt", "", ui)
	if err != nil {
		return err
	}
	bundle, ok := k.(*libkb.PgpKeyBundle)
	if ok {
		s.keys = append(s.keys, (*openpgp.Entity)(bundle))
	}
	return nil
}

// scan finds the user on the api server for the key id.  Then it
// uses PGPKeyfinder to find the public pgp keys for the user,
// identifying/tracking along the way.
func (s *ScanKeys) scan(id uint64) (openpgp.EntityList, error) {
	// lookup the user on the api server by the key id.
	username, uid, err := s.apiLookup(id)
	if err != nil {
		return nil, err
	}
	G.Log.Info("key id %d => (%s, %s)", id, username, uid)

	// use PGPKeyfinder engine to get the pgp keys for the user
	arg := &PGPKeyfinderArg{Users: []string{"uid://" + uid}}
	if s.opts != nil {
		arg.TrackOptions = *s.opts
	}
	ctx := &Context{SecretUI: s.secui, IdentifyUI: s.idui}
	eng := NewPGPKeyfinder(arg)
	if err := RunEngine(eng, ctx); err != nil {
		return nil, err
	}
	uplus := eng.UsersPlusKeys()
	if len(uplus) != 1 {
		G.Log.Warning("error getting user plus pgp key from %s", username)
		return nil, err
	}

	// convert the bundles to an openpgp entity list
	// (which implements the openpgp.KeyRing interface)
	var list openpgp.EntityList
	for _, k := range uplus[0].Keys {
		list = append(list, (*openpgp.Entity)(k))
	}
	return list, nil

}

// apiLookup gets the username and uid from the api server for the
// key id.
func (s *ScanKeys) apiLookup(id uint64) (username, uid string, err error) {
	var data struct {
		Username string
		UID      string
	}

	// lookup key on api server
	args := libkb.ApiArg{
		Endpoint: "key/basics",
		Args: libkb.HttpArgs{
			"pgp_key_id": libkb.UHex{Val: id},
		},
	}
	if err = G.API.GetDecode(args, &data); err != nil {
		return "", "", err
	}
	G.Log.Debug("response data: %+v", data)
	return data.Username, data.UID, nil
}
