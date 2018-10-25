
export default class AdminListOrgsCtrl {
  /** @ngInject */
  constructor($scope, backendSrv, navModelSrv) {
    $scope.init = () => {
      $scope.navModel = navModelSrv.getNav('cfg', 'admin', 'global-orgs', 1);
      $scope.getOrgs();
    };

    $scope.getOrgs = () => {
      backendSrv.get('/api/orgs').then(orgs => {
        $scope.orgs = orgs;
      });
    };

    $scope.deleteOrg = org => {
      $scope.appEvent('confirm-modal', {
        title: 'Löschen',
        text: 'Wollen Sie die Organisation ' + org.name + ' löschen?',
        text2: 'Alle Dashboards von dieser Organisation werden entfernt!',
        icon: 'fa-trash',
        yesText: 'Löschen',
        onConfirm: () => {
          backendSrv.delete('/api/orgs/' + org.id).then(() => {
            $scope.getOrgs();
          });
        },
      });
    };

    $scope.init();
  }
}

