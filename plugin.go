package minrpc

import "net/http"

// Plugin 是一个插件接口
type Plugin interface {
    Init() error
    Handle(next http.Handler) http.Handler
}

// PluginManager 是一个插件管理器
type PluginManager struct {
    plugins []Plugin
}

// NewPluginManager 创建一个新的插件管理器
func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make([]Plugin, 0),
    }
}

// RegisterPlugin 注册一个插件
func (m *PluginManager) RegisterPlugin(p Plugin) error {
    if err := p.Init(); err != nil {
        return err
    }
    m.plugins = append(m.plugins, p)
    return nil
}

// Then 将插件应用到处理器
func (m *PluginManager) Then(h http.Handler) http.Handler {
    for i := len(m.plugins) - 1; i >= 0; i-- {
        h = m.plugins[i].Handle(h)
    }
    return h
}