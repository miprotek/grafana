import React from 'react';
import { inject, observer } from 'mobx-react';
import { toJS } from 'mobx';
import PageHeader from 'app/core/components/PageHeader/PageHeader';
import IContainerProps from 'app/containers/IContainerProps';
import { getSnapshot } from 'mobx-state-tree';
import appEvents from 'app/core/app_events';

@inject('nav', 'folder', 'view')
@observer
export class FolderSettings extends React.Component<IContainerProps, any> {
  formSnapshot: any;

  constructor(props) {
    super(props);
    this.loadStore();
  }

  loadStore() {
    const { nav, folder, view } = this.props;

    return folder.load(view.routeParams.get('uid') as string).then(res => {
      this.formSnapshot = getSnapshot(folder);
      view.updatePathAndQuery(`${res.url}/settings`, {}, {});

      return nav.initFolderNav(toJS(folder.folder), 'manage-folder-settings');
    });
  }

  onTitleChange(evt) {
    this.props.folder.setTitle(this.getFormSnapshot().folder.title, evt.target.value);
  }

  getFormSnapshot() {
    if (!this.formSnapshot) {
      this.formSnapshot = getSnapshot(this.props.folder);
    }

    return this.formSnapshot;
  }

  save(evt) {
    if (evt) {
      evt.stopPropagation();
      evt.preventDefault();
    }

    const { nav, folder, view } = this.props;

    folder
      .saveFolder({ overwrite: false })
      .then(newUrl => {
        view.updatePathAndQuery(newUrl, {}, {});

        appEvents.emit('dashboard-saved');
        appEvents.emit('alert-success', ['Ordner gespeichert']);
      })
      .then(() => {
        return nav.initFolderNav(toJS(folder.folder), 'manage-folder-settings');
      })
      .catch(this.handleSaveFolderError.bind(this));
  }

  delete(evt) {
    if (evt) {
      evt.stopPropagation();
      evt.preventDefault();
    }

    const { folder, view } = this.props;
    const title = folder.folder.title;

    appEvents.emit('confirm-modal', {
      title: 'Löschen',
      text: `Wollen Sie diesen Ordner und alle enthaltenen Dashboards löschen?`,
      icon: 'fa-trash',
      yesText: 'Löschen',
      onConfirm: () => {
        return folder.deleteFolder().then(() => {
          appEvents.emit('alert-success', ['Ordner gelöscht', `${title} wurde gelöscht`]);
          view.updatePathAndQuery('dashboards', '', '');
        });
      },
    });
  }

  handleSaveFolderError(err) {
    if (err.data && err.data.status === 'version-mismatch') {
      err.isHandled = true;

      const { nav, folder, view } = this.props;

      appEvents.emit('confirm-modal', {
        title: 'Konflikt',
        text: 'Jemand anderes hat diesen Ordner aktualisiert.',
        text2: 'Wollen Sie diesen Ordner trotzdem speichern?',
        yesText: 'Speichern & Überschreiben',
        icon: 'fa-warning',
        onConfirm: () => {
          folder
            .saveFolder({ overwrite: true })
            .then(newUrl => {
              view.updatePathAndQuery(newUrl, {}, {});

              appEvents.emit('dashboard-saved');
              appEvents.emit('alert-success', ['Ordner gespeichert']);
            })
            .then(() => {
              return nav.initFolderNav(toJS(folder.folder), 'manage-folder-settings');
            });
        },
      });
    }
  }

  render() {
    const { nav, folder } = this.props;

    if (!folder.folder || !nav.main) {
      return <h2>Loading</h2>;
    }

    return (
      <div>
        <PageHeader model={nav as any} />
        <div className="page-container page-body">
          <h2 className="page-sub-heading">Ordnereinstellungen</h2>

          <div className="section gf-form-group">
            <form name="folderSettingsForm" onSubmit={this.save.bind(this)}>
              <div className="gf-form">
                <label className="gf-form-label width-7">Name</label>
                <input
                  type="text"
                  className="gf-form-input width-30"
                  value={folder.folder.title}
                  onChange={this.onTitleChange.bind(this)}
                />
              </div>
              <div className="gf-form-button-row">
                <button
                  type="submit"
                  className="btn btn-success"
                  disabled={!folder.folder.canSave || !folder.folder.hasChanged}
                >
                  <i className="fa fa-save" /> Speichern
                </button>
                <button className="btn btn-danger" onClick={this.delete.bind(this)} disabled={!folder.folder.canSave}>
                  <i className="fa fa-trash" /> Löschen
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
}
