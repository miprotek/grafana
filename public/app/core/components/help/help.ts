import coreModule from '../../core_module';
import appEvents from 'app/core/app_events';

export class HelpCtrl {
  tabIndex: any;
  shortcuts: any;

  /** @ngInject */
  constructor() {
    this.tabIndex = 0;
    this.shortcuts = {
      Global: [
        { keys: ['g', 'h'], description: 'Gehe zu Home Dashboard' },
        { keys: ['g', 'p'], description: 'Gehe zum Profil' },
        { keys: ['s', 'o'], description: 'Suche öffnen' },
        { keys: ['s', 's'], description: 'Öffne Suche mit Favorisierten Filter' },
        { keys: ['s', 't'], description: 'Öffnet die Suche in der Ansicht mit Tags' },
        { keys: ['esc'], description: 'Beendet Bearbeitungs- und Einstellungsansichten' },
      ],
      Dashboard: [
        { keys: ['mod+s'], description: 'Dashboard speichern' },
        { keys: ['d', 'r'], description: 'Alle Panels aktualisieren' },
        { keys: ['d', 's'], description: 'Dashboard Einstellungen' },
        { keys: ['d', 'v'], description: 'Anischtsmodus an- oder ausschalten' },
        { keys: ['d', 'k'], description: 'Kioskmodus umschalten (versteckt obere Navigationsleiste)' },
        { keys: ['d', 'E'], description: 'Alle Zeilen erweitern' },
        { keys: ['d', 'C'], description: 'Alle Zeilen reduzieren' },
        { keys: ['mod+o'], description: 'Geteiltes Diagramm-Fadenkreuz umschalten' },
      ],
      'Focused Panel': [
        { keys: ['e'], description: 'Ansichtsfenster bearbeiten' },
        { keys: ['v'], description: 'Ansicht des Vollbildmodus umschalten' },
        { keys: ['p', 's'], description: 'Öffnet das Freigabe Modal' },
        { keys: ['p', 'd'], description: 'Panel duplizieren' },
        { keys: ['p', 'r'], description: 'Panel entfernen' },
      ],
      'Time Range': [
        { keys: ['t', 'z'], description: 'Zeitbereich verkleinern' },
        {
          keys: ['t', '<i class="fa fa-long-arrow-left"></i>'],
          description: 'Zeitbereich zurück schieben',
        },
        {
          keys: ['t', '<i class="fa fa-long-arrow-right"></i>'],
          description: 'Zeitbereich vor schieben',
        },
      ],
    };
  }

  dismiss() {
    appEvents.emit('hide-modal');
  }
}

export function helpModal() {
  return {
    restrict: 'E',
    templateUrl: 'public/app/core/components/help/help.html',
    controller: HelpCtrl,
    bindToController: true,
    transclude: true,
    controllerAs: 'ctrl',
    scope: {},
  };
}

coreModule.directive('helpModal', helpModal);
