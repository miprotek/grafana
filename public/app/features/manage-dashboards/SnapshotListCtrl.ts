import _ from 'lodash';

export class SnapshotListCtrl {
  navModel: any;
  snapshots: any;

  /** @ngInject */
  constructor(private $rootScope, private backendSrv, navModelSrv) {
    this.navModel = navModelSrv.getNav('dashboards', 'snapshots', 0);
    this.backendSrv.get('/api/dashboard/snapshots').then(result => {
      this.snapshots = result;
    });
  }

  removeSnapshotConfirmed(snapshot) {
    _.remove(this.snapshots, { key: snapshot.key });
    this.backendSrv.delete('/api/snapshots/' + snapshot.key).then(
      () => {},
      () => {
        this.snapshots.push(snapshot);
      }
    );
  }

  removeSnapshot(snapshot) {
    this.$rootScope.appEvent('confirm-modal', {
      title: 'Löschen',
      text: 'Sind Sie sicher, dass Sie den Snapshot ' + snapshot.name + ' löschen wollen?',
      yesText: 'Löschen',
      icon: 'fa-trash',
      onConfirm: () => {
        this.removeSnapshotConfirmed(snapshot);
      },
    });
  }
}
