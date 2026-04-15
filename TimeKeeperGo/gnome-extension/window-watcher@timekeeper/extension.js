import Gio from 'gi://Gio';

const IFACE_XML = `
<node>
  <interface name="io.timekeeper.WindowWatcher">
    <method name="GetActiveWindow">
      <arg type="s" direction="out" name="title"/>
    </method>
  </interface>
</node>`;

class WindowWatcherImpl {
    GetActiveWindow() {
        const win = global.display.focus_window;
        return win ? (win.get_title() ?? '') : '';
    }
}

export default class WindowWatcherExtension {
    enable() {
        this._impl = new WindowWatcherImpl();
        this._exported = Gio.DBusExportedObject.wrapJSObject(IFACE_XML, this._impl);
        this._exported.export(Gio.DBus.session, '/io/timekeeper/WindowWatcher');
        this._nameId = Gio.bus_own_name_on_connection(
            Gio.DBus.session,
            'io.timekeeper.WindowWatcher',
            Gio.BusNameOwnerFlags.NONE,
            null, null
        );
    }

    disable() {
        Gio.bus_unown_name(this._nameId);
        this._exported.unexport();
        this._exported = null;
        this._impl = null;
    }
}
