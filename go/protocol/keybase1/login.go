// Auto-generated by avdl-compiler v1.3.21 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/keybase1/login.avdl

package keybase1

import (
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type ConfiguredAccount struct {
	Username        string `codec:"username" json:"username"`
	HasStoredSecret bool   `codec:"hasStoredSecret" json:"hasStoredSecret"`
}

func (o ConfiguredAccount) DeepCopy() ConfiguredAccount {
	return ConfiguredAccount{
		Username:        o.Username,
		HasStoredSecret: o.HasStoredSecret,
	}
}

type GetConfiguredAccountsArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type LoginArg struct {
	SessionID       int        `codec:"sessionID" json:"sessionID"`
	DeviceType      string     `codec:"deviceType" json:"deviceType"`
	UsernameOrEmail string     `codec:"usernameOrEmail" json:"usernameOrEmail"`
	ClientType      ClientType `codec:"clientType" json:"clientType"`
}

type LoginProvisionedDeviceArg struct {
	SessionID          int    `codec:"sessionID" json:"sessionID"`
	Username           string `codec:"username" json:"username"`
	NoPassphrasePrompt bool   `codec:"noPassphrasePrompt" json:"noPassphrasePrompt"`
}

type LoginWithPaperKeyArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type ClearStoredSecretArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Username  string `codec:"username" json:"username"`
}

type LogoutArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type DeprovisionArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Username  string `codec:"username" json:"username"`
	DoRevoke  bool   `codec:"doRevoke" json:"doRevoke"`
}

type RecoverAccountFromEmailAddressArg struct {
	Email string `codec:"email" json:"email"`
}

type PaperKeyArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type PaperKeySubmitArg struct {
	SessionID   int    `codec:"sessionID" json:"sessionID"`
	PaperPhrase string `codec:"paperPhrase" json:"paperPhrase"`
}

type UnlockArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type UnlockWithPassphraseArg struct {
	SessionID  int    `codec:"sessionID" json:"sessionID"`
	Passphrase string `codec:"passphrase" json:"passphrase"`
}

type PGPProvisionArg struct {
	SessionID  int    `codec:"sessionID" json:"sessionID"`
	Username   string `codec:"username" json:"username"`
	Passphrase string `codec:"passphrase" json:"passphrase"`
	DeviceName string `codec:"deviceName" json:"deviceName"`
}

type AccountDeleteArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type LoginInterface interface {
	// Returns an array of information about accounts configured on the local
	// machine. Currently configured accounts are defined as those that have stored
	// secrets, but this definition may be expanded in the future.
	GetConfiguredAccounts(context.Context, int) ([]ConfiguredAccount, error)
	// Performs login.  deviceType should be libkb.DeviceTypeDesktop
	// or libkb.DeviceTypeMobile.  usernameOrEmail is optional.
	// If the current device isn't provisioned, this function will
	// provision it.
	//
	// Note that if usernameOrEmail is an email address, only provisioning
	// will be attempted.  If the device is already provisioned, login
	// via email address does not work.
	Login(context.Context, LoginArg) error
	// Login a user only if the user is on a provisioned device.  Username is optional.
	// If noPassphrasePrompt is set, then only a stored secret will be used to unlock
	// the device keys.
	LoginProvisionedDevice(context.Context, LoginProvisionedDeviceArg) error
	// Login and unlock by
	// - trying unlocked device keys if available
	// - prompting for a paper key and using that
	LoginWithPaperKey(context.Context, int) error
	// Removes any existing stored secret for the given username.
	// loginWithStoredSecret(_, username) will fail after this is called.
	ClearStoredSecret(context.Context, ClearStoredSecretArg) error
	Logout(context.Context, int) error
	Deprovision(context.Context, DeprovisionArg) error
	RecoverAccountFromEmailAddress(context.Context, string) error
	// PaperKey generates paper backup keys for restoring an account.
	// It calls login_ui.displayPaperKeyPhrase with the phrase.
	PaperKey(context.Context, int) error
	// paperKeySubmit checks that paperPhrase is a valid paper key
	// for the logged in user, caches the keys, and sends a notification.
	PaperKeySubmit(context.Context, PaperKeySubmitArg) error
	// Unlock restores access to local key store by priming passphrase stream cache.
	Unlock(context.Context, int) error
	UnlockWithPassphrase(context.Context, UnlockWithPassphraseArg) error
	// pgpProvision is for devel/testing to provision a device via pgp using CLI
	// with no user interaction.
	PGPProvision(context.Context, PGPProvisionArg) error
	// accountDelete is for devel/testing to delete the current user's account.
	AccountDelete(context.Context, int) error
}

func LoginProtocol(i LoginInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "keybase.1.login",
		Methods: map[string]rpc.ServeHandlerDescription{
			"getConfiguredAccounts": {
				MakeArg: func() interface{} {
					ret := make([]GetConfiguredAccountsArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetConfiguredAccountsArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetConfiguredAccountsArg)(nil), args)
						return
					}
					ret, err = i.GetConfiguredAccounts(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"login": {
				MakeArg: func() interface{} {
					ret := make([]LoginArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoginArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoginArg)(nil), args)
						return
					}
					err = i.Login(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loginProvisionedDevice": {
				MakeArg: func() interface{} {
					ret := make([]LoginProvisionedDeviceArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoginProvisionedDeviceArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoginProvisionedDeviceArg)(nil), args)
						return
					}
					err = i.LoginProvisionedDevice(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loginWithPaperKey": {
				MakeArg: func() interface{} {
					ret := make([]LoginWithPaperKeyArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoginWithPaperKeyArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoginWithPaperKeyArg)(nil), args)
						return
					}
					err = i.LoginWithPaperKey(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"clearStoredSecret": {
				MakeArg: func() interface{} {
					ret := make([]ClearStoredSecretArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ClearStoredSecretArg)
					if !ok {
						err = rpc.NewTypeError((*[]ClearStoredSecretArg)(nil), args)
						return
					}
					err = i.ClearStoredSecret(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"logout": {
				MakeArg: func() interface{} {
					ret := make([]LogoutArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LogoutArg)
					if !ok {
						err = rpc.NewTypeError((*[]LogoutArg)(nil), args)
						return
					}
					err = i.Logout(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"deprovision": {
				MakeArg: func() interface{} {
					ret := make([]DeprovisionArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]DeprovisionArg)
					if !ok {
						err = rpc.NewTypeError((*[]DeprovisionArg)(nil), args)
						return
					}
					err = i.Deprovision(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"recoverAccountFromEmailAddress": {
				MakeArg: func() interface{} {
					ret := make([]RecoverAccountFromEmailAddressArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]RecoverAccountFromEmailAddressArg)
					if !ok {
						err = rpc.NewTypeError((*[]RecoverAccountFromEmailAddressArg)(nil), args)
						return
					}
					err = i.RecoverAccountFromEmailAddress(ctx, (*typedArgs)[0].Email)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"paperKey": {
				MakeArg: func() interface{} {
					ret := make([]PaperKeyArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]PaperKeyArg)
					if !ok {
						err = rpc.NewTypeError((*[]PaperKeyArg)(nil), args)
						return
					}
					err = i.PaperKey(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"paperKeySubmit": {
				MakeArg: func() interface{} {
					ret := make([]PaperKeySubmitArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]PaperKeySubmitArg)
					if !ok {
						err = rpc.NewTypeError((*[]PaperKeySubmitArg)(nil), args)
						return
					}
					err = i.PaperKeySubmit(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"unlock": {
				MakeArg: func() interface{} {
					ret := make([]UnlockArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]UnlockArg)
					if !ok {
						err = rpc.NewTypeError((*[]UnlockArg)(nil), args)
						return
					}
					err = i.Unlock(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"unlockWithPassphrase": {
				MakeArg: func() interface{} {
					ret := make([]UnlockWithPassphraseArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]UnlockWithPassphraseArg)
					if !ok {
						err = rpc.NewTypeError((*[]UnlockWithPassphraseArg)(nil), args)
						return
					}
					err = i.UnlockWithPassphrase(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"pgpProvision": {
				MakeArg: func() interface{} {
					ret := make([]PGPProvisionArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]PGPProvisionArg)
					if !ok {
						err = rpc.NewTypeError((*[]PGPProvisionArg)(nil), args)
						return
					}
					err = i.PGPProvision(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"accountDelete": {
				MakeArg: func() interface{} {
					ret := make([]AccountDeleteArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]AccountDeleteArg)
					if !ok {
						err = rpc.NewTypeError((*[]AccountDeleteArg)(nil), args)
						return
					}
					err = i.AccountDelete(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type LoginClient struct {
	Cli rpc.GenericClient
}

// Returns an array of information about accounts configured on the local
// machine. Currently configured accounts are defined as those that have stored
// secrets, but this definition may be expanded in the future.
func (c LoginClient) GetConfiguredAccounts(ctx context.Context, sessionID int) (res []ConfiguredAccount, err error) {
	__arg := GetConfiguredAccountsArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.login.getConfiguredAccounts", []interface{}{__arg}, &res)
	return
}

// Performs login.  deviceType should be libkb.DeviceTypeDesktop
// or libkb.DeviceTypeMobile.  usernameOrEmail is optional.
// If the current device isn't provisioned, this function will
// provision it.
//
// Note that if usernameOrEmail is an email address, only provisioning
// will be attempted.  If the device is already provisioned, login
// via email address does not work.
func (c LoginClient) Login(ctx context.Context, __arg LoginArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.login", []interface{}{__arg}, nil)
	return
}

// Login a user only if the user is on a provisioned device.  Username is optional.
// If noPassphrasePrompt is set, then only a stored secret will be used to unlock
// the device keys.
func (c LoginClient) LoginProvisionedDevice(ctx context.Context, __arg LoginProvisionedDeviceArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.loginProvisionedDevice", []interface{}{__arg}, nil)
	return
}

// Login and unlock by
// - trying unlocked device keys if available
// - prompting for a paper key and using that
func (c LoginClient) LoginWithPaperKey(ctx context.Context, sessionID int) (err error) {
	__arg := LoginWithPaperKeyArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.login.loginWithPaperKey", []interface{}{__arg}, nil)
	return
}

// Removes any existing stored secret for the given username.
// loginWithStoredSecret(_, username) will fail after this is called.
func (c LoginClient) ClearStoredSecret(ctx context.Context, __arg ClearStoredSecretArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.clearStoredSecret", []interface{}{__arg}, nil)
	return
}

func (c LoginClient) Logout(ctx context.Context, sessionID int) (err error) {
	__arg := LogoutArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.login.logout", []interface{}{__arg}, nil)
	return
}

func (c LoginClient) Deprovision(ctx context.Context, __arg DeprovisionArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.deprovision", []interface{}{__arg}, nil)
	return
}

func (c LoginClient) RecoverAccountFromEmailAddress(ctx context.Context, email string) (err error) {
	__arg := RecoverAccountFromEmailAddressArg{Email: email}
	err = c.Cli.Call(ctx, "keybase.1.login.recoverAccountFromEmailAddress", []interface{}{__arg}, nil)
	return
}

// PaperKey generates paper backup keys for restoring an account.
// It calls login_ui.displayPaperKeyPhrase with the phrase.
func (c LoginClient) PaperKey(ctx context.Context, sessionID int) (err error) {
	__arg := PaperKeyArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.login.paperKey", []interface{}{__arg}, nil)
	return
}

// paperKeySubmit checks that paperPhrase is a valid paper key
// for the logged in user, caches the keys, and sends a notification.
func (c LoginClient) PaperKeySubmit(ctx context.Context, __arg PaperKeySubmitArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.paperKeySubmit", []interface{}{__arg}, nil)
	return
}

// Unlock restores access to local key store by priming passphrase stream cache.
func (c LoginClient) Unlock(ctx context.Context, sessionID int) (err error) {
	__arg := UnlockArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.login.unlock", []interface{}{__arg}, nil)
	return
}

func (c LoginClient) UnlockWithPassphrase(ctx context.Context, __arg UnlockWithPassphraseArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.unlockWithPassphrase", []interface{}{__arg}, nil)
	return
}

// pgpProvision is for devel/testing to provision a device via pgp using CLI
// with no user interaction.
func (c LoginClient) PGPProvision(ctx context.Context, __arg PGPProvisionArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.login.pgpProvision", []interface{}{__arg}, nil)
	return
}

// accountDelete is for devel/testing to delete the current user's account.
func (c LoginClient) AccountDelete(ctx context.Context, sessionID int) (err error) {
	__arg := AccountDeleteArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.login.accountDelete", []interface{}{__arg}, nil)
	return
}
