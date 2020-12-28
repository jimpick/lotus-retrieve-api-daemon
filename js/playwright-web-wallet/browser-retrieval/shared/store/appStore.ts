import { makeAutoObservable } from 'mobx'
import { Node } from 'shared/Node'

import { AlertsStore } from './alertsStore'
import { DealsStore } from './dealsStore'
import { LogsStore } from './logsStore'
import { OffersStore } from './offersStore'
import { OperationsStore } from './operationsStore'
import { OptionsStore } from './optionsStore'
import { QueriesStore } from './queriesStore'
import { RecentCIDStore } from './recentCIDStore'
import { SettingsStore } from './settingsStore'
import { UploadStore } from './uploadStore'
import { PCHStore } from './pchStore'

export class AppStore {
  alertsStore: AlertsStore
  dealsStore: DealsStore
  logsStore: LogsStore
  offersStore: OffersStore
  operationsStore: OperationsStore
  optionsStore: OptionsStore
  queriesStore: QueriesStore
  recentCIDStore: RecentCIDStore
  settingsStore: SettingsStore
  uploadStore: UploadStore
  pchStore: PCHStore

  node: Node

  connected = false

  connectedPeers = []

  constructor() {
    this.alertsStore = new AlertsStore(this)
    this.dealsStore = new DealsStore(this)
    this.logsStore = new LogsStore(this)
    this.offersStore = new OffersStore(this)
    this.operationsStore = new OperationsStore(this)
    this.optionsStore = new OptionsStore(this)
    this.queriesStore = new QueriesStore(this)
    this.recentCIDStore = new RecentCIDStore(this)
    this.settingsStore = new SettingsStore(this)
    this.uploadStore = new UploadStore(this)
    this.pchStore = new PCHStore(this)

    makeAutoObservable(this)
  }

  async connect() {
    const { optionsStore } = this

    if (optionsStore.wallet || optionsStore.privateKey) {
      try {
        await this.tryDisconnect()

        this.node = await Node.create(true)
        this.connected = true

        const mimBalance = await this.node.lotus.hasMinBalance()

        if (!mimBalance) {
          this.alertsStore.create({
            message: 'Your wallet does not have minimum FIL required!',
            type: 'warning',
          })
        }
      } catch (error) {
        this.connected = false

        if (error === 'Error: `Invalid Key Length`') {
          this.logsStore.logDebug(`start node failed: ${error}`)
          console.error(error)
        } else {
          console.error(error)
          this.logsStore.logError(`start node failed: ${error.message}`)
        }
      }
    }
  }

  private async tryDisconnect() {
    try {
      await this.node?.stop()
      this.connected = false
    } catch (err) {
      this.logsStore.logDebug(`stop node failed: ${err.message}`)
    }
  }

  query({ cid, minerID }) {
    this.node.query(cid, minerID)
  }

  setConnectedPeers(connectedPeers) {
    this.connectedPeers = connectedPeers
  }

  downloadFile(msg: { cid: any; offer?: any }) {
    this.node.downloadFile(msg)
  }

  deleteFile({ cid }) {
    this.node.deleteFile(cid)
  }

  stopAutomationCode() {
    this.node.stopLoop()
  }

  runAutomationCode() {
    this.node.runAutomationCode()
  }
}

export const appStore = new AppStore()
