import { types, getEnv, flow } from 'mobx-state-tree';
import { ServerStat } from './ServerStat';

export const ServerStatsStore = types
  .model('ServerStatsStore', {
    stats: types.array(ServerStat),
    error: types.optional(types.string, ''),
  })
  .actions(self => ({
    load: flow(function* load() {
      const backendSrv = getEnv(self).backendSrv;
      const res = yield backendSrv.get('/api/admin/stats');
      self.stats.clear();
      self.stats.push(ServerStat.create({ name: 'Insgesamt Dashboards', value: res.dashboards }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der Nutzer', value: res.users }));
      self.stats.push(ServerStat.create({ name: 'Aktive Benutzer (Innerhalb der letzten 30 Tage)', value: res.activeUsers }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der Organisationen', value: res.orgs }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der Playlisten', value: res.playlists }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der Snapshots', value: res.snapshots }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der Dashboard tags', value: res.tags }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der markierten Dashboards', value: res.stars }));
      self.stats.push(ServerStat.create({ name: 'Gesamtzahl der Meldungen', value: res.alerts }));
    }),
  }));
