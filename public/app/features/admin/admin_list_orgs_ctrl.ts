import angular from 'angular';

export class AdminListOrgsCtrl {
  /** @ngInject */
  constructor($scope, backendSrv, navModelSrv) {
    $scope.init = function() {
      $scope.navModel = navModelSrv.getNav('cfg', 'admin', 'global-orgs', 1);
      $scope.getOrgs();
    };

    $scope.getOrgs = function() {
      backendSrv.get('/api/orgs').then(function(orgs) {
        $scope.orgs = orgs;
      });
    };

    $scope.deleteOrg = function(org) {
      $scope.appEvent('confirm-modal', {
        title: 'Löschen',
        text: 'Wollen Sie die Organisation ' + org.name + ' löschen?',
        text2: 'Alle Dashboards von dieser Organisation werden entfernt!',
        icon: 'fa-trash',
        yesText: 'Löschen',
        onConfirm: function() {
          backendSrv.delete('/api/orgs/' + org.id).then(function() {
            $scope.getOrgs();
          });
        },
      });
    };

    $scope.init();
  }
}

angular.module('grafana.controllers').controller('AdminListOrgsCtrl', AdminListOrgsCtrl);
