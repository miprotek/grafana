import coreModule from '../core_module';

export class ResetPasswordCtrl {
  /** @ngInject */
  constructor($scope, contextSrv, backendSrv, $location) {
    contextSrv.sidemenu = false;
    $scope.formModel = {};
    $scope.mode = 'send';

    const params = $location.search();
    if (params.code) {
      $scope.mode = 'reset';
      $scope.formModel.code = params.code;
    }

    $scope.navModel = {
      main: {
        icon: 'gicon gicon-branding',
        text: 'Passwort zurücksetzen',
        subTitle: 'Setzen Sie ihr Grafana Passwort zurück',
        breadcrumbs: [{ title: 'Login', url: 'login' }],
      },
    };

    $scope.sendResetEmail = () => {
      if (!$scope.sendResetForm.$valid) {
        return;
      }
      backendSrv.post('/api/user/password/send-reset-email', $scope.formModel).then(() => {
        $scope.mode = 'email-sent';
      });
    };

    $scope.submitReset = () => {
      if (!$scope.resetForm.$valid) {
        return;
      }

      if ($scope.formModel.newPassword !== $scope.formModel.confirmPassword) {
        $scope.appEvent('alert-warning', ['Neue Passwörter stimmen nicht überein', '']);
        return;
      }

      backendSrv.post('/api/user/password/reset', $scope.formModel).then(() => {
        $location.path('login');
      });
    };
  }
}

coreModule.controller('ResetPasswordCtrl', ResetPasswordCtrl);
