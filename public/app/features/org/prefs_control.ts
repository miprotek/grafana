import config from 'app/core/config';
import coreModule from 'app/core/core_module';

export class PrefsControlCtrl {
  prefs: any;
  oldTheme: any;
  prefsForm: any;
  mode: string;

  timezones: any = [
    { value: '', text: 'Default' },
    { value: 'browser', text: 'Lokale Browser Zeit' },
    { value: 'utc', text: 'UTC' },
  ];
  themes: any = [{ value: '', text: 'Standard' }, { value: 'dark', text: 'Dunkel' }, { value: 'light', text: 'Hell' }];

  /** @ngInject **/
  constructor(private backendSrv, private $location) {}

  $onInit() {
    return this.backendSrv.get(`/api/${this.mode}/preferences`).then(prefs => {
      this.prefs = prefs;
      this.oldTheme = prefs.theme;
    });
  }

  updatePrefs() {
    if (!this.prefsForm.$valid) {
      return;
    }

    var cmd = {
      theme: this.prefs.theme,
      timezone: this.prefs.timezone,
      homeDashboardId: this.prefs.homeDashboardId,
    };

    this.backendSrv.put(`/api/${this.mode}/preferences`, cmd).then(() => {
      window.location.href = config.appSubUrl + this.$location.path();
    });
  }
}

var template = `
<form name="ctrl.prefsForm" class="section gf-form-group">
  <h3 class="page-heading">Einstellungen</h3>

  <div class="gf-form">
    <span class="gf-form-label width-11">UI Theme</span>
    <div class="gf-form-select-wrapper max-width-20">
      <select class="gf-form-input" ng-model="ctrl.prefs.theme" ng-options="f.value as f.text for f in ctrl.themes"></select>
    </div>
  </div>

  <div class="gf-form">
    <span class="gf-form-label width-11">
      Home Dashboard
      <info-popover mode="right-normal">
        Finden Sie nicht das gewünschte Dashboard? Markieren Sie es zuerst, dann sollte es in diesem Auswahlfeld erscheinen.
      </info-popover>
    </span>
    <dashboard-selector class="gf-form-select-wrapper max-width-20" model="ctrl.prefs.homeDashboardId">
    </dashboard-selector>
  </div>

  <div class="gf-form">
    <label class="gf-form-label width-11">Zeitzone</label>
    <div class="gf-form-select-wrapper max-width-20">
      <select class="gf-form-input" ng-model="ctrl.prefs.timezone" ng-options="f.value as f.text for f in ctrl.timezones"></select>
    </div>
  </div>

  <div class="gf-form-button-row">
    <button type="submit" class="btn btn-success" ng-click="ctrl.updatePrefs()">Sichern</button>
  </div>
</form>
`;

export function prefsControlDirective() {
  return {
    restrict: 'E',
    controller: PrefsControlCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    template: template,
    scope: {
      mode: '@',
    },
  };
}

coreModule.directive('prefsControl', prefsControlDirective);
