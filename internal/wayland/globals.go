package wayland

type Globals struct {
	registry                             *WlRegistry
	wlShm                                *WlShm
	zwpLinuxDmabufV1                     *ZwpLinuxDmabufV1
	wlCompositor                         *WlCompositor
	wlSubcompositor                      *WlSubcompositor
	wlDataDeviceManager                  *WlDataDeviceManager
	zxdgOutputManagerV1                  *ZxdgOutputManagerV1
	zwpIdleInhibitManagerV1              *ZwpIdleInhibitManagerV1
	xdgWmBase                            *XdgWmBase
	zwpTabletManagerV2                   *ZwpTabletManagerV2
	zxdgDecorationManagerV1              *ZxdgDecorationManagerV1
	zwpRelativePointerManagerV1          *ZwpRelativePointerManagerV1
	zwpPointerConstraintsV1              *ZwpPointerConstraintsV1
	wpPresentation                       *WpPresentation
	zwpTextInputManagerV3                *ZwpTextInputManagerV3
	zwpPrimarySelectionDeviceManagerV1   *ZwpPrimarySelectionDeviceManagerV1
	wpViewporter                         *WpViewporter
	zwpKeyboardShortcutsInhibitManagerV1 *ZwpKeyboardShortcutsInhibitManagerV1
	wlSeat                               *WlSeat
	wlOutput                             *WlOutput

	globals map[string]WlRegistryGlobalEvent
	conn    *Display
}

func (g *Globals) registerGlobal(event *WlRegistryGlobalEvent) {
	g.globals[event.Interface] = *event
}

func (g *Globals) unregisterGlobal(event *WlRegistryGlobalRemoveEvent) {
	for intf, global := range g.globals {
		if global.Name == event.Name {
			delete(g.globals, intf)
			return
		}
	}
}

func (g *Globals) Registry() *WlRegistry {
	if g.registry != nil {
		return g.registry
	}
	registry, err := g.conn.display.GetRegistry(g.conn)
	if err != nil {
		panic(err)
	}
	g.registry = registry
	if err := g.conn.Sync(); err != nil {
		panic(err)
	}
	return registry
}

func (g *Globals) WlShm() *WlShm {
	registry := g.Registry()
	if global, ok := g.globals[WlShmDescriptor.Name]; ok {
		if g.wlShm != nil {
			return g.wlShm
		}
		id, err := registry.Bind(g.conn, global.Name, WlShmDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WlShm{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		return proxy
	}
	return nil
}

func (g *Globals) ZwpLinuxDmabufV1() *ZwpLinuxDmabufV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpLinuxDmabufV1Descriptor.Name]; ok {
		if g.zwpLinuxDmabufV1 != nil {
			return g.zwpLinuxDmabufV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpLinuxDmabufV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpLinuxDmabufV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpLinuxDmabufV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WlCompositor() *WlCompositor {
	registry := g.Registry()
	if global, ok := g.globals[WlCompositorDescriptor.Name]; ok {
		if g.wlCompositor != nil {
			return g.wlCompositor
		}
		id, err := registry.Bind(g.conn, global.Name, WlCompositorDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WlCompositor{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wlCompositor = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WlSubcompositor() *WlSubcompositor {
	registry := g.Registry()
	if global, ok := g.globals[WlSubcompositorDescriptor.Name]; ok {
		if g.wlSubcompositor != nil {
			return g.wlSubcompositor
		}
		id, err := registry.Bind(g.conn, global.Name, WlSubcompositorDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WlSubcompositor{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wlSubcompositor = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WlDataDeviceManager() *WlDataDeviceManager {
	registry := g.Registry()
	if global, ok := g.globals[WlDataDeviceManagerDescriptor.Name]; ok {
		if g.wlDataDeviceManager != nil {
			return g.wlDataDeviceManager
		}
		id, err := registry.Bind(g.conn, global.Name, WlDataDeviceManagerDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WlDataDeviceManager{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wlDataDeviceManager = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZxdgOutputManagerV1() *ZxdgOutputManagerV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZxdgOutputManagerV1Descriptor.Name]; ok {
		if g.zxdgOutputManagerV1 != nil {
			return g.zxdgOutputManagerV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZxdgOutputManagerV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZxdgOutputManagerV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zxdgOutputManagerV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpIdleInhibitManagerV1() *ZwpIdleInhibitManagerV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpIdleInhibitManagerV1Descriptor.Name]; ok {
		if g.zwpIdleInhibitManagerV1 != nil {
			return g.zwpIdleInhibitManagerV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpIdleInhibitManagerV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpIdleInhibitManagerV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpIdleInhibitManagerV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) XdgWmBase() *XdgWmBase {
	registry := g.Registry()
	if global, ok := g.globals[XdgWmBaseDescriptor.Name]; ok {
		if g.xdgWmBase != nil {
			return g.xdgWmBase
		}
		id, err := registry.Bind(g.conn, global.Name, XdgWmBaseDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &XdgWmBase{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.xdgWmBase = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpTabletManagerV2() *ZwpTabletManagerV2 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpTabletManagerV2Descriptor.Name]; ok {
		if g.zwpTabletManagerV2 != nil {
			return g.zwpTabletManagerV2
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpTabletManagerV2Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpTabletManagerV2{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpTabletManagerV2 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZxdgDecorationManagerV1() *ZxdgDecorationManagerV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZxdgDecorationManagerV1Descriptor.Name]; ok {
		if g.zxdgDecorationManagerV1 != nil {
			return g.zxdgDecorationManagerV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZxdgDecorationManagerV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZxdgDecorationManagerV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zxdgDecorationManagerV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpRelativePointerManagerV1() *ZwpRelativePointerManagerV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpRelativePointerManagerV1Descriptor.Name]; ok {
		if g.zwpRelativePointerManagerV1 != nil {
			return g.zwpRelativePointerManagerV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpRelativePointerManagerV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpRelativePointerManagerV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpRelativePointerManagerV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpPointerConstraintsV1() *ZwpPointerConstraintsV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpPointerConstraintsV1Descriptor.Name]; ok {
		if g.zwpPointerConstraintsV1 != nil {
			return g.zwpPointerConstraintsV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpPointerConstraintsV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpPointerConstraintsV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpPointerConstraintsV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WpPresentation() *WpPresentation {
	registry := g.Registry()
	if global, ok := g.globals[WpPresentationDescriptor.Name]; ok {
		if g.wpPresentation != nil {
			return g.wpPresentation
		}
		id, err := registry.Bind(g.conn, global.Name, WpPresentationDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WpPresentation{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wpPresentation = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpTextInputManagerV3() *ZwpTextInputManagerV3 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpTextInputManagerV3Descriptor.Name]; ok {
		if g.zwpTextInputManagerV3 != nil {
			return g.zwpTextInputManagerV3
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpTextInputManagerV3Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpTextInputManagerV3{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpTextInputManagerV3 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpPrimarySelectionDeviceManagerV1() *ZwpPrimarySelectionDeviceManagerV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpPrimarySelectionDeviceManagerV1Descriptor.Name]; ok {
		if g.zwpPrimarySelectionDeviceManagerV1 != nil {
			return g.zwpPrimarySelectionDeviceManagerV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpPrimarySelectionDeviceManagerV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpPrimarySelectionDeviceManagerV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpPrimarySelectionDeviceManagerV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WpViewporter() *WpViewporter {
	registry := g.Registry()
	if global, ok := g.globals[WpViewporterDescriptor.Name]; ok {
		if g.wpViewporter != nil {
			return g.wpViewporter
		}
		id, err := registry.Bind(g.conn, global.Name, WpViewporterDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WpViewporter{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wpViewporter = proxy
		return proxy
	}
	return nil
}

func (g *Globals) ZwpKeyboardShortcutsInhibitManagerV1() *ZwpKeyboardShortcutsInhibitManagerV1 {
	registry := g.Registry()
	if global, ok := g.globals[ZwpKeyboardShortcutsInhibitManagerV1Descriptor.Name]; ok {
		if g.zwpKeyboardShortcutsInhibitManagerV1 != nil {
			return g.zwpKeyboardShortcutsInhibitManagerV1
		}
		id, err := registry.Bind(g.conn, global.Name, ZwpKeyboardShortcutsInhibitManagerV1Descriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &ZwpKeyboardShortcutsInhibitManagerV1{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.zwpKeyboardShortcutsInhibitManagerV1 = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WlSeat() *WlSeat {
	registry := g.Registry()
	if global, ok := g.globals[WlSeatDescriptor.Name]; ok {
		if g.wlSeat != nil {
			return g.wlSeat
		}
		id, err := registry.Bind(g.conn, global.Name, WlSeatDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WlSeat{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wlSeat = proxy
		return proxy
	}
	return nil
}

func (g *Globals) WlOutput() *WlOutput {
	registry := g.Registry()
	if global, ok := g.globals[WlOutputDescriptor.Name]; ok {
		if g.wlOutput != nil {
			return g.wlOutput
		}
		id, err := registry.Bind(g.conn, global.Name, WlOutputDescriptor.Name, global.Version)
		if err != nil {
			panic(err)
		}
		proxy := &WlOutput{id: id, version: global.Version}
		g.conn.RegisterProxy(proxy)
		g.wlOutput = proxy
		return proxy
	}
	return nil
}
