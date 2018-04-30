import coreModule from '../../core/core_module';
import _ from 'lodash';

export class DataSourcesCtrl {
  datasources: any;
  unfiltered: any;
  navModel: any;
  searchQuery: string;

  /** @ngInject */
  constructor(private $scope, private backendSrv, private datasourceSrv, private navModelSrv) {
    this.navModel = this.navModelSrv.getNav('cfg', 'datasources', 0);
    backendSrv.get('/api/datasources').then(result => {
      this.datasources = result;
      this.unfiltered = result;
    });
  }

  onQueryUpdated() {
    let regex = new RegExp(this.searchQuery, 'ig');
    this.datasources = _.filter(this.unfiltered, item => {
      return regex.test(item.name) || regex.test(item.type);
    });
  }

  removeDataSourceConfirmed(ds) {
    this.backendSrv
      .delete('/api/datasources/' + ds.id)
      .then(
        () => {
          this.$scope.appEvent('alert-success', ['Datenquelle gelöscht', '']);
        },
        () => {
          this.$scope.appEvent('alert-error', ['Datenquelle konnte nicht gelöscht werden', '']);
        }
      )
      .then(() => {
        this.backendSrv.get('/api/datasources').then(result => {
          this.datasources = result;
        });
        this.backendSrv.get('/api/frontend/settings').then(settings => {
          this.datasourceSrv.init(settings.datasources);
        });
      });
  }

  removeDataSource(ds) {
    this.$scope.appEvent('confirm-modal', {
      title: 'Löschen',
      text: 'Sind Sie sicher, dass Sie die Datenquelle ' + ds.name + ' löschen wollen?',
      yesText: 'Löschen',
      icon: 'fa-trash',
      onConfirm: () => {
        this.removeDataSourceConfirmed(ds);
      },
    });
  }
}

coreModule.controller('DataSourcesCtrl', DataSourcesCtrl);
