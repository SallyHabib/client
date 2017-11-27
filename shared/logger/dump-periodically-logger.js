// @flow
import type {Logger, LogFn, LogLevel, LogLineWithLevel} from './types'
import {requestIdleCallback} from '../util/idle-callback'

type FileWriterFn = (lines: Array<LogLineWithLevel>) => Promise<void>

// Dumps the inner logger periodically, everything else is forwarded
// At most every `periodInMs` seconds. May be a little less because
// requestIdleCallback is used.
class DumpPeriodicallyLogger implements Logger {
  _innerLogger: Logger
  _periodInMs: number
  _fileWriterFn: FileWriterFn
  _lastTimeoutId: ?number
  _levelPrefix: LogLevel
  _ok: boolean = true
  log: LogFn
  dump: (levelPrefix: LogLevel) => Promise<Array<LogLineWithLevel>> // Should return an ordered array of log lines (ordered by timestamp)

  constructor(innerLogger: Logger, periodInMs: number, fileWriterFn: FileWriterFn, levelPrefix: LogLevel) {
    this._innerLogger = innerLogger
    this._periodInMs = periodInMs
    this._fileWriterFn = fileWriterFn
    this._levelPrefix = levelPrefix
    this.log = innerLogger.log
    this.dump = innerLogger.dump
    this._periodicallyDump()
  }

  _periodicallyDump = () => {
    if (this._ok) {
      return this._innerLogger
        .dump(this._levelPrefix)
        .then(this._fileWriterFn)
        .then(() => {
          this._lastTimeoutId = setTimeout(
            () =>
              requestIdleCallback(
                deadline => {
                  this._periodicallyDump()
                },
                {timeout: this._periodInMs}
              ),
            this._periodInMs
          )
        })
        .catch(e => {
          console.error('dump-periodically failed', e)
          this._ok = false
        })
    }

    return Promise.reject(new Error('Not ok'))
  }

  flush() {
    this._ok = true
    this._lastTimeoutId && clearTimeout(this._lastTimeoutId)
    return this._innerLogger.flush().then(this._periodicallyDump)
  }
}

export default DumpPeriodicallyLogger