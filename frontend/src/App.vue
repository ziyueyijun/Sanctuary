<!--
  Sanctuary - 主组件
  
  功能：
  - 服务管理：启动/停止 Nginx、MySQL、PHP-CGI
  - 站点管理：添加、编辑、删除站点配置
  - 系统设置：查看端口配置信息
  
  作者：子曰亦君
  QQ：15593838
-->
<script setup>
// ==================== 导入依赖 ====================
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { IconDesktop, IconApps, IconSettings, IconPlus } from '@arco-design/web-vue/es/icon'
// 导入 Go 后端方法
import {
  GetServices, StartService, StopService, StartAll, StopAll,
  RestartAll, OpenWWWFolder, OpenBrowser, GetWWWPath, GetDefaultSiteRoot, SelectDirectory,
  GetSites, UpdateSites, OpenSiteInBrowser, OpenSiteFolder,
  GetSettings, UpdateSettings, OpenConfigFile, OpenLogFile, OpenLogsFolder
} from '../wailsjs/go/main/App'

// ==================== 服务官方 Logo SVG ====================
const serviceIcons = {
  // Nginx 官方绿色 Logo
  Nginx: `<svg viewBox="0 0 32 32"><path fill="#009639" d="M16 2L3 9v14l13 7 13-7V9L16 2zm5.9 18.6c0 .4-.2.8-.5 1.1-.3.3-.7.5-1.2.5-.6 0-1.1-.3-1.4-.8l-4.4-7v6.7c0 .6-.4 1.1-1.1 1.1h-.1c-.6 0-1.1-.5-1.1-1.1V11.4c0-.4.2-.8.5-1.1.3-.3.7-.5 1.2-.5.6 0 1.1.3 1.4.8l4.4 7v-6.7c0-.6.5-1.1 1.1-1.1h.1c.6 0 1.1.5 1.1 1.1v9.7z"/></svg>`,
  // MySQL 官方海豚 Logo
  MySQL: `<svg viewBox="0 0 32 32"><path fill="#00758F" d="M27.7 24.8c-1.2-.1-2.1 0-2.9.4-.2.1-.6.1-.6.4 0 .2.2.3.4.5.3.4.7.9 1.1 1.2.4.3.9.6 1.3.9.8.5 1.7.7 2.4 1.1.4.3.8.5 1.2.8.2.1.3.3.5.3v-.1c-.1-.2-.2-.3-.3-.5l-.6-.6c-.6-.8-1.4-1.5-2.3-2-.7-.4-2.3-1-2.6-1.8l-.1-.1c.5-.1 1.1-.2 1.6-.4.8-.2 1.5-.2 2.3-.4.4-.1.7-.2 1.1-.3v-.2c-.4-.4-.7-.9-1.2-1.3-1.2-1-2.5-2-3.9-2.9-.7-.5-1.7-.8-2.4-1.2-.3-.2-.7-.2-.9-.5-.4-.5-.6-1.1-.9-1.7-1.1-2.2-2.2-4.7-3.1-7.1-.7-1.6-1.1-3.3-2-4.7-4.1-6.5-8.5-10.4-15.4-14.3-1.4-.8-3.2-1.1-5-1.5-.3 0-.7-.1-1.1-.1-.5-.2-.9-.8-1.3-1.1C3.1.9 1.3-.3 2 1.1c-.3.9.5 1.8.8 2.3.6.8 1.2 1.5 1.6 2.4.3.6.3 1.2.5 1.8.5 1.4.9 2.9 1.6 4.2.3.7.7 1.3 1.1 1.9.2.4.6.6.7 1.1-.4.6-.4 1.4-.6 2.1-.9 3.2-.6 7.2.8 9.6.4.7 1.5 2.2 2.9 1.6 1.2-.5 1-2 1.3-3.4.1-.3 0-.5.2-.7v.1c.4.7.7 1.5 1.1 2.2.8 1.3 2.3 2.7 3.5 3.6.6.5 1.2 1.3 2 1.6v-.1h-.1c-.2-.2-.4-.4-.6-.6-.5-.5-1.1-1.1-1.5-1.6-1.2-1.5-2.3-3.2-3.2-4.9-.5-.9-.9-1.8-1.2-2.7-.2-.4-.2-.9-.4-1.2-.4.6-.9 1.1-1.2 1.8-.5 1.1-.5 2.5-.7 3.9-.1 0-.1 0-.2-.1-.9-.4-1.2-1.4-1.5-2.3-.8-2.4-.9-6.2 0-8.9.2-.7 1.2-2.9.8-3.5-.2-.6-.9-1-1.2-1.5-.4-.6-.9-1.4-1.2-2.1-.8-1.9-1.2-4-2.1-5.9-.4-.9-1.1-1.8-1.7-2.6-.6-.9-1.3-1.5-1.8-2.6-.2-.4-.4-1-.2-1.4.1-.3.2-.4.5-.5.5-.4 1.8.1 2.3.3 1.4.5 2.6 1 3.7 1.7.5.4 1.1.9 1.7 1.1h.8c1.2.3 2.5.1 3.6.5 2 .7 3.8 1.7 5.4 2.8 4.9 3.5 8.9 8.4 11.7 14.3.4.9.6 1.8 1 2.8.7 1.9 1.6 3.9 2.3 5.8.7 1.8 1.4 3.7 2.5 5.2.5.8 2.6 1.2 3.6 1.7.7.3 1.8.7 2.4 1.1.7.4 1.3.8 1.9 1.2"/></svg>`,
  // PHP 官方紫色 Logo
  'PHP-CGI': `<svg viewBox="0 0 32 32"><ellipse cx="16" cy="16" rx="14" ry="8" fill="#777BB4"/><text x="16" y="19" text-anchor="middle" fill="#fff" font-size="8" font-weight="bold" font-family="Arial">php</text></svg>`
}

const getServiceIcon = (name) => {
  return serviceIcons[name] || serviceIcons['Nginx']
}

// ==================== 响应式数据 ====================
const activeMenu = ref(['services'])    // 当前激活的菜单
// 预设默认服务数据，避免首屏等待
const services = ref([
  { name: 'Nginx', version: '1.26.2', status: 'stopped', port: '80' },
  { name: 'MySQL', version: '8.4.7', status: 'stopped', port: '3306' },
  { name: 'PHP-CGI', version: '8.1.33', status: 'stopped', port: '9001' }
])
const wwwPath = ref('')                 // 网站根目录
// 按钮独立加载状态
const startAllLoading = ref(false)
const stopAllLoading = ref(false)
const restartAllLoading = ref(false)
const loadingServices = ref({})         // 按服务名存储加载状态
// 站点管理
const sites = ref([])
const siteDialogVisible = ref(false)
const isEditSite = ref(false)
const editingSiteId = ref(0)
const siteForm = ref({ name: '', domain: '', port: 80, root: '', rewriteRules: [] })
const siteTabActiveKey = ref('basic')  // 站点编辑对话框标签页
const rewriteRulesText = ref('')       // URL 重写规则文本
// URL 重写类型选项
const rewriteTypes = ref([
  { label: 'Rewrite', value: 'rewrite' },
  { label: 'Redirect', value: 'redirect' },
  { label: 'Proxy', value: 'proxy' }
])
// 系统设置
const settings = ref({
  autoStart: false,
  nginxPort: 80,
  phpPort: 9001,
  mysqlPort: 3306
})
const settingsSaving = ref(false)
// 定时器
let timer = null

// ==================== 表格列配置 ====================
const siteColumns = [
  { title: '站点名称', dataIndex: 'name', width: 150 },
  { title: '域名', dataIndex: 'domain', width: 150, slotName: 'domain' },
  { title: '端口', dataIndex: 'port', width: 80 },
  { title: '根目录', dataIndex: 'root', ellipsis: true, tooltip: true, slotName: 'root' },
  { title: '状态', dataIndex: 'enabled', width: 100, slotName: 'status' },
  { title: '操作', width: 150, slotName: 'actions' }
]

const refreshServices = async () => {
  try {
    services.value = await GetServices()
  } catch (e) {
    console.error(e)
  }
}

const startService = async (name) => {
  loadingServices.value[name] = true
  try {
    await StartService(name)
    await new Promise(r => setTimeout(r, 500))
    await refreshServices()
    Message.success(`${name} 启动成功`)
  } catch (e) {
    Message.error(`启动 ${name} 失败: ${e}`)
  }
  loadingServices.value[name] = false
}

const stopService = async (name) => {
  loadingServices.value[name] = true
  try {
    await StopService(name)
    await new Promise(r => setTimeout(r, 500))
    await refreshServices()
    Message.success(`${name} 已停止`)
  } catch (e) {
    Message.error(`停止 ${name} 失败: ${e}`)
  }
  loadingServices.value[name] = false
}

const startAll = async () => {
  startAllLoading.value = true
  try {
    await StartAll()
    await new Promise(r => setTimeout(r, 1000))
    await refreshServices()
    Message.success('所有服务已启动')
  } catch (e) {
    Message.error(`启动失败: ${e}`)
  }
  startAllLoading.value = false
}

const stopAll = async () => {
  stopAllLoading.value = true
  try {
    await StopAll()
    await new Promise(r => setTimeout(r, 500))
    await refreshServices()
    Message.success('所有服务已停止')
  } catch (e) {
    Message.error(`停止失败: ${e}`)
  }
  stopAllLoading.value = false
}

const restartAll = async () => {
  restartAllLoading.value = true
  try {
    await RestartAll()
    await new Promise(r => setTimeout(r, 1000))
    await refreshServices()
    Message.success('所有服务已重启')
  } catch (e) {
    Message.error(`重启失败: ${e}`)
  }
  restartAllLoading.value = false
}

const openSiteInBrowser = (site) => {
  OpenSiteInBrowser(site.domain, site.port)
}

const openSiteFolder = (site) => {
  OpenSiteFolder(site.root)
}

const showAddSite = async () => {
  isEditSite.value = false
  siteForm.value = { name: '', domain: '', port: 80, root: '', rewriteRules: [] }
  rewriteRulesText.value = ''
  siteTabActiveKey.value = 'basic'
  siteDialogVisible.value = true
}

const showEditSite = (site) => {
  isEditSite.value = true
  editingSiteId.value = site.id
  siteForm.value = { 
    ...site,
    rewriteRules: site.rewriteRules ? site.rewriteRules.map(rule => ({ ...rule })) : []
  }
  
  // 将现有的重写规则转换为文本格式
  rewriteRulesText.value = siteForm.value.rewriteRules.map(rule => {
    if (rule.type === 'rewrite') {
      return `rewrite ^${rule.pattern}$ ${rule.destination} last;`
    } else if (rule.type === 'redirect') {
      return `rewrite ^${rule.pattern}$ ${rule.destination} permanent;`
    } else if (rule.type === 'proxy') {
      return `location ~ ^${rule.pattern}$ {
    proxy_pass ${rule.destination};
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}`
    }
    return ''
  }).join('\n')
  
  siteTabActiveKey.value = 'basic'
  siteDialogVisible.value = true
}

const selectSiteRoot = async () => {
  try {
    const dir = await SelectDirectory()
    if (dir) {
      siteForm.value.root = dir
    }
  } catch (e) {
    console.error(e)
  }
}

const saveSite = async () => {
  if (!siteForm.value.name || !siteForm.value.domain || !siteForm.value.root) {
    Message.warning('请填写完整信息')
    return
  }
  
  // 如果在 URL 重写标签页且有文本内容，先解析规则
  if (siteTabActiveKey.value === 'rewrite' && rewriteRulesText.value) {
    parseRewriteRules()
  }
  
  if (isEditSite.value) {
    const index = sites.value.findIndex(s => s.id === editingSiteId.value)
    if (index !== -1) {
      sites.value[index] = { ...siteForm.value, id: editingSiteId.value, enabled: sites.value[index].enabled }
    }
  } else {
    sites.value.push({
      id: Date.now(),
      ...siteForm.value,
      enabled: true
    })
  }
  siteDialogVisible.value = false
  
  try {
    await UpdateSites(sites.value)
    Message.success(isEditSite.value ? '站点已更新，Nginx 已重启' : '站点已添加，Nginx 已重启')
  } catch (e) {
    Message.error(`保存失败: ${e}`)
  }
}

const deleteSite = (site) => {
  Modal.warning({
    title: '提示',
    content: `确定删除站点 "${site.name}" 吗？`,
    okText: '确定',
    cancelText: '取消',
    hideCancel: false,
    onOk: async () => {
      sites.value = sites.value.filter(s => s.id !== site.id)
      try {
        await UpdateSites(sites.value)
        Message.success('站点已删除，Nginx 已重启')
      } catch (e) {
        Message.error(`删除失败: ${e}`)
      }
    }
  })
}

watch(() => siteForm.value.domain, async (newDomain) => {
  if (!isEditSite.value) {
    if (newDomain) {
      const www = await GetWWWPath()
      siteForm.value.root = `${www}/${newDomain}`
    } else {
      siteForm.value.root = ''
    }
  }
})

// ==================== 设置相关方法 ====================
const loadSettings = async () => {
  try {
    const data = await GetSettings()
    settings.value = data
  } catch (e) {
    console.error('加载设置失败:', e)
  }
}

const saveSettings = async () => {
  settingsSaving.value = true
  try {
    await UpdateSettings(settings.value)
    Message.success('设置已保存')
  } catch (e) {
    Message.error(`保存失败: ${e}`)
  }
  settingsSaving.value = false
}

onMounted(async () => {
  wwwPath.value = await GetWWWPath()
  sites.value = await GetSites()
  await loadSettings()
  await refreshServices()
  timer = setInterval(refreshServices, 3000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

// URL 重写规则相关函数
const parseRewriteRules = () => {
  const rules = []
  const lines = rewriteRulesText.value.split('\n')
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i].trim()
    if (!line || line.startsWith('#')) continue
    
    // 解析 rewrite 规则
    const rewriteMatch = line.match(/^rewrite\s+\^(.*?)\$\s+(.*?)(?:\s+(last|permanent|redirect));?$/i)
    if (rewriteMatch) {
      const [, pattern, destination, type] = rewriteMatch
      rules.push({
        id: Date.now() + i,
        pattern: pattern,
        destination: destination,
        type: type === 'permanent' || type === 'redirect' ? 'redirect' : 'rewrite',
        enabled: true
      })
      continue
    }
    
    // 解析 location proxy 规则 (简化处理)
    const proxyMatch = line.match(/location\s+~\s+\^(.*?)\$\s*\{.*?proxy_pass\s+(.*?);/i)
    if (proxyMatch) {
      const [, pattern, destination] = proxyMatch
      rules.push({
        id: Date.now() + i,
        pattern: pattern,
        destination: destination,
        type: 'proxy',
        enabled: true
      })
      continue
    }
  }
  
  siteForm.value.rewriteRules = rules
  Message.success(`成功解析 ${rules.length} 条规则`)
}

const clearRewriteRules = () => {
  rewriteRulesText.value = ''
  siteForm.value.rewriteRules = []
}

const removeRewriteRule = (rule) => {
  const index = siteForm.value.rewriteRules.findIndex(r => r.id === rule.id)
  if (index !== -1) {
    siteForm.value.rewriteRules.splice(index, 1)
  }
}

</script>

<template>
  <a-layout class="app-container">
    <!-- 左侧菜单 -->
    <a-layout-sider :width="200" class="aside">
      <div class="logo">
        <h2>Sanctuary</h2>
      </div>
      <a-menu
        v-model:selected-keys="activeMenu"
        :style="{ width: '100%' }"
      >
        <a-menu-item key="services">
          <template #icon><icon-desktop /></template>
          服务管理
        </a-menu-item>
        <a-menu-item key="sites">
          <template #icon><icon-apps /></template>
          站点管理
        </a-menu-item>
        <a-menu-item key="settings">
          <template #icon><icon-settings /></template>
          系统设置
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <!-- 右侧内容 -->
    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <div class="app-logo">
            <svg width="24" height="24" viewBox="0 0 256 256" style="margin-right: 8px; vertical-align: middle;">
              <defs>
                <linearGradient id="headerBgGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#4F46E5"/>
                  <stop offset="100%" stop-color="#7C3AED"/>
                </linearGradient>
              </defs>
              <circle cx="128" cy="128" r="120" fill="url(#headerBgGradient)"/>
              <path d="M90 70 L166 70 Q180 70, 180 85 L180 100 Q180 115, 166 115 L110 115 L110 140 L160 140 Q175 140, 175 155 L175 170 Q175 185, 160 185 L90 185 Q75 185, 75 170 L75 155 Q75 140, 90 140 L145 140 L145 115 L95 115 Q80 115, 80 100 L80 85 Q80 70, 95 70 Z" fill="white"/>
            </svg>
            <span v-if="activeMenu[0] === 'services'">服务管理</span>
            <span v-else-if="activeMenu[0] === 'sites'">站点管理</span>
            <span v-else>系统设置</span>
          </div>
        </div>
      </a-layout-header>

      <a-layout-content class="main">
        <!-- 服务管理 -->
        <div v-if="activeMenu[0] === 'services'">
          <div class="toolbar">
            <a-button type="primary" status="success" @click="startAll" :loading="startAllLoading">启动全部</a-button>
            <a-button type="primary" status="danger" @click="stopAll" :loading="stopAllLoading">停止全部</a-button>
            <a-button type="primary" status="warning" @click="restartAll" :loading="restartAllLoading">重启全部</a-button>
          </div>

          <div class="service-cards">
            <a-card v-for="service in services" :key="service.name" class="service-card" hoverable>
              <div class="service-header">
                <div class="service-icon" :class="service.status === 'running' ? 'running' : 'stopped'">
                  <span class="custom-icon" v-html="getServiceIcon(service.name)"></span>
                </div>
                <div class="service-info">
                  <h3>{{ service.name }}</h3>
                  <p>{{ service.version }}</p>
                </div>
                <a-tag :color="service.status === 'running' ? 'green' : 'red'" size="small">
                  {{ service.status === 'running' ? '运行中' : '已停止' }}
                </a-tag>
              </div>
              <div class="service-detail">
                <span>端口：{{ service.port }}</span>
              </div>
              <div class="service-actions">
                <a-button
                  v-if="service.status !== 'running'"
                  type="primary"
                  status="success"
                  @click="startService(service.name)"
                  :loading="loadingServices[service.name]"
                  long
                >启动</a-button>
                <a-button
                  v-else
                  type="primary"
                  status="danger"
                  @click="stopService(service.name)"
                  :loading="loadingServices[service.name]"
                  long
                >停止</a-button>
              </div>
            </a-card>
          </div>
          
          <!-- 开发者信息 -->
          <div class="footer-info">
            <span>开发者：子曰亦君</span>
            <span class="divider">|</span>
            <span>QQ：15593838</span>
          </div>
        </div>

        <!-- 站点管理 -->
        <div v-else-if="activeMenu[0] === 'sites'">
          <div class="toolbar">
            <a-button type="primary" @click="showAddSite">
              <template #icon><icon-plus /></template>
              添加站点
            </a-button>
          </div>

          <a-table :columns="siteColumns" :data="sites" :bordered="true" stripe>
            <template #domain="{ record }">
              <span class="domain-link" @click="openSiteInBrowser(record)">{{ record.domain }}</span>
            </template>
            <template #root="{ record }">
              <span class="root-path" @click="openSiteFolder(record)">{{ record.root }}</span>
            </template>
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">
                {{ record.enabled ? '已启用' : '已禁用' }}
              </a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button size="small" type="primary" @click="showEditSite(record)">编辑</a-button>
                <a-button size="small" type="primary" status="danger" @click="deleteSite(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </div>

        <!-- 系统设置 -->
        <div v-else-if="activeMenu[0] === 'settings'" class="settings-page">
          <a-tabs default-active-key="general" animation>
            <!-- 常规设置 -->
            <a-tab-pane key="general" title="常规设置">
              <a-card class="settings-card">
                <a-form :model="settings" layout="horizontal" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }">
                  <a-form-item label="开机启动">
                    <a-switch v-model="settings.autoStart" />
                    <span class="form-tip">系统启动时自动运行 Sanctuary</span>
                  </a-form-item>
                </a-form>
              </a-card>
            </a-tab-pane>

            <!-- 端口设置 -->
            <a-tab-pane key="ports" title="端口设置">
              <a-card class="settings-card">
                <a-form :model="settings" layout="horizontal" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }">
                  <a-form-item label="Nginx 端口">
                    <a-input-number v-model="settings.nginxPort" :min="1" :max="65535" style="width: 150px" />
                  </a-form-item>
                  <a-form-item label="PHP 端口">
                    <a-input-number v-model="settings.phpPort" :min="1" :max="65535" style="width: 150px" />
                  </a-form-item>
                  <a-form-item label="MySQL 端口">
                    <a-input-number v-model="settings.mysqlPort" :min="1" :max="65535" style="width: 150px" />
                  </a-form-item>
                </a-form>
              </a-card>
            </a-tab-pane>

          </a-tabs>

          <!-- 保存按钮 -->
          <div class="settings-actions">
            <a-button type="primary" :loading="settingsSaving" @click="saveSettings">保存设置</a-button>
          </div>
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>

  <!-- 添加/编辑站点对话框 -->
  <a-modal
    v-model:visible="siteDialogVisible"
    :title="isEditSite ? '编辑站点' : '添加站点'"
    :width="700"
    @ok="saveSite"
    @cancel="siteDialogVisible = false"
  >
    <a-tabs v-model:activeKey="siteTabActiveKey">
      <a-tab-pane key="basic" title="基本信息">
        <a-form :model="siteForm" layout="horizontal" :label-col-props="{ span: 5 }" :wrapper-col-props="{ span: 19 }">
          <a-form-item label="站点名称">
            <a-input v-model="siteForm.name" placeholder="例如：公司网站、个人博客" />
          </a-form-item>
          <a-form-item label="域名">
            <a-input v-model="siteForm.domain" placeholder="例如：www.php.net" />
          </a-form-item>
          <a-form-item label="端口">
            <a-input-number v-model="siteForm.port" :min="1" :max="65535" />
          </a-form-item>
          <a-form-item label="根目录">
            <a-input-group>
              <a-input v-model="siteForm.root" style="flex: 1" />
              <a-button @click="selectSiteRoot">选择</a-button>
            </a-input-group>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      
      <a-tab-pane key="rewrite" title="URL重写">
        <div class="rewrite-rules-container">
          <a-alert type="info" style="margin-bottom: 16px;">
            <template #icon><icon-info-circle-fill /></template>
            您可以直接粘贴 Nginx 格式的重写规则，每行一条规则。支持的格式：
            <br><code>rewrite ^/pattern$ /destination last;</code>
            <br><code>rewrite ^/pattern$ /destination permanent;</code>
          </a-alert>
          
          <a-textarea 
            v-model="rewriteRulesText" 
            placeholder="在此粘贴您的 URL 重写规则，例如：
rewrite ^/api/(.*)$ /api.php?path=$1 last;
rewrite ^/admin$ /admin.php permanent;"
            :auto-size="{ minRows: 6, maxRows: 12 }"
            style="margin-bottom: 16px;"
          />
          
          <div class="rewrite-actions">
            <a-button @click="parseRewriteRules">解析规则</a-button>
            <a-button style="margin-left: 12px;" @click="clearRewriteRules">清空规则</a-button>
          </div>
          
          <a-divider />
          
          <div style="margin-bottom: 16px;">
            <strong>已解析的规则 ({{ siteForm.rewriteRules.length }} 条)：</strong>
          </div>
          
          <a-table 
            :data="siteForm.rewriteRules" 
            :pagination="false"
            :bordered="false"
            row-key="id"
            v-if="siteForm.rewriteRules.length > 0"
          >
            <template #columns>
              <a-table-column title="启用" :width="60">
                <template #cell="{ record }">
                  <a-checkbox v-model="record.enabled" />
                </template>
              </a-table-column>
              
              <a-table-column title="匹配模式" :width="150">
                <template #cell="{ record }">
                  <a-input v-model="record.pattern" placeholder="/api/(.*)" />
                </template>
              </a-table-column>
              
              <a-table-column title="类型" :width="100">
                <template #cell="{ record }">
                  <a-select v-model="record.type" :options="rewriteTypes" style="width: 100%;" />
                </template>
              </a-table-column>
              
              <a-table-column title="目标地址">
                <template #cell="{ record }">
                  <a-input v-model="record.destination" placeholder="/index.php?route=$1" />
                </template>
              </a-table-column>
              
              <a-table-column title="操作" :width="80">
                <template #cell="{ record }">
                  <a-button type="text" status="danger" @click="removeRewriteRule(record)">
                    <template #icon><icon-delete /></template>
                  </a-button>
                </template>
              </a-table-column>
            </template>
          </a-table>
          
          <div v-else style="text-align: center; padding: 24px; color: #999;">
            暂无解析的规则，请在上方文本框中粘贴规则并点击"解析规则"
          </div>
        </div>
      </a-tab-pane>
    </a-tabs>
  </a-modal>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: 'Microsoft YaHei', 'PingFang SC', -apple-system, BlinkMacSystemFont, sans-serif;
}

.app-container {
  height: 100%;
}

/* 侧边栏 - 深紫渐变背景 */
.aside {
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 50%, #1a1a2e 100%) !important;
}

.aside .arco-layout-sider-children {
  display: flex;
  flex-direction: column;
}

/* Logo - 渐变紫 */
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.logo h2 {
  color: #fff;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 2px;
  text-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

/* 菜单样式 */
.aside .arco-menu {
  background: transparent !important;
  flex: 1;
  border-right: none !important;
}

.aside .arco-menu-inner {
  background: transparent !important;
}

.aside .arco-menu-item {
  margin: 4px 12px;
  border-radius: 8px;
  color: rgba(255,255,255,0.7) !important;
  background: transparent !important;
  transition: all 0.3s ease;
}

.aside .arco-menu-item:hover {
  background: rgba(102, 126, 234, 0.2) !important;
  color: #fff !important;
}

.aside .arco-menu-item:hover .arco-icon {
  color: #fff !important;
}

/* 激活菜单项 - 渐变紫 */
.aside .arco-menu-selected {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: #fff !important;
  font-weight: 500;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.aside .arco-menu-item .arco-icon {
  color: rgba(255,255,255,0.7);
}

.aside .arco-menu-selected .arco-icon {
  color: #fff;
}

/* 页头 */
.header {
  background: #fff !important;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 64px !important;
}

.header-left {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

/* 主内容区 */
.main {
  background: #f5f3ff;
  padding: 24px;
  overflow-y: auto;
}

.toolbar {
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
}

/* 服务卡片 */
.service-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.service-card {
  border-radius: 12px !important;
  transition: all 0.3s ease;
  border: 1px solid #e9e4f5 !important;
}

.service-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(102, 126, 234, 0.15) !important;
  border-color: #c4b5fd !important;
}

.service-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 16px;
}

.service-icon {
  width: 52px;
  height: 52px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border: 2px solid #e5e7eb;
  transition: all 0.3s ease;
}

.custom-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.custom-icon svg {
  width: 100%;
  height: 100%;
}

/* 运行中 - 绿色边框 */
.service-icon.running {
  background: #f0fdf4;
  border-color: #22c55e;
  box-shadow: 0 0 12px rgba(34, 197, 94, 0.3);
}

/* 已停止 - 灰色边框 */
.service-icon.stopped {
  background: #f9fafb;
  border-color: #d1d5db;
  opacity: 0.7;
}

.service-info {
  flex: 1;
}

.service-info h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.service-info p {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
}

.service-detail {
  padding: 12px;
  background: #faf5ff;
  border-radius: 8px;
  margin-bottom: 14px;
  color: #4b5563;
  font-size: 13px;
}

.service-actions {
  display: flex;
}

/* 域名链接 - 紫色 */
.domain-link {
  color: #7c3aed;
  cursor: pointer;
  transition: all 0.2s;
}

.domain-link:hover {
  color: #8b5cf6;
  text-decoration: underline;
}

/* 根目录 */
.root-path {
  cursor: pointer;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.2s;
}

.root-path:hover {
  color: #7c3aed;
}

/* 开发者信息 */
.footer-info {
  margin-top: 30px;
  padding: 16px 0;
  text-align: center;
  color: #909399;
  font-size: 13px;
  border-top: 1px solid #e9e4f5;
}

.footer-info .divider {
  margin: 0 12px;
  color: #dcdfe6;
}

/* Arco Design 组件多彩主题覆盖 */

/* 主色按钮 - 渐变紫 */
.arco-btn-primary:not(.arco-btn-status-success):not(.arco-btn-status-danger):not(.arco-btn-status-warning) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  border: none !important;
}

.arco-btn-primary:not(.arco-btn-status-success):not(.arco-btn-status-danger):not(.arco-btn-status-warning):hover {
  background: linear-gradient(135deg, #7c8ff0 0%, #8b5cb8 100%) !important;
}

/* 成功按钮 - 绿色 */
.arco-btn-status-success {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
  border: none !important;
}

.arco-btn-status-success:hover {
  background: linear-gradient(135deg, #34d399 0%, #10b981 100%) !important;
}

/* 危险按钮 - 红色 */
.arco-btn-status-danger {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%) !important;
  border: none !important;
}

.arco-btn-status-danger:hover {
  background: linear-gradient(135deg, #f87171 0%, #ef4444 100%) !important;
}

/* 警告按钮 - 橙色 */
.arco-btn-status-warning {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%) !important;
  border: none !important;
}

.arco-btn-status-warning:hover {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%) !important;
}

/* 次要按钮 - 边框紫色 */
.arco-btn-secondary {
  border-color: #8b5cf6 !important;
  color: #7c3aed !important;
}

.arco-btn-secondary:hover {
  background: #f5f3ff !important;
  border-color: #7c3aed !important;
}

/* 表格悬停行 - 淡紫色 */
.arco-table-tr:hover .arco-table-td {
  background: #faf5ff !important;
}

/* 表格头部 */
.arco-table-th {
  background: #f8f7ff !important;
  color: #4c1d95 !important;
  font-weight: 600 !important;
}

/* 标签绿色 */
.arco-tag-green {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%) !important;
  color: #047857 !important;
  border: none !important;
}

/* 标签红色 */
.arco-tag-red {
  background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%) !important;
  color: #b91c1c !important;
  border: none !important;
}

/* 标签灰色 */
.arco-tag-gray {
  background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%) !important;
  color: #4b5563 !important;
  border: none !important;
}

/* 弹窗标题 */
.arco-modal-header {
  border-bottom: 1px solid #e9e4f5;
}

/* 表单输入框聚焦 */
.arco-input-wrapper:focus-within {
  border-color: #8b5cf6 !important;
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2) !important;
}

.arco-input-number:focus-within {
  border-color: #8b5cf6 !important;
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2) !important;
}

/* ==================== 设置页面样式 ==================== */
.settings-page {
  max-width: 800px;
}

.settings-card {
  border-radius: 12px !important;
  border: 1px solid #e9e4f5 !important;
}

.settings-card .arco-card-body {
  padding: 20px;
}

.settings-card .arco-card-header {
  border-bottom: 1px solid #f0ebfa;
}

.settings-card .arco-card-header-title {
  font-weight: 600;
  color: #374151;
}

.form-tip {
  margin-left: 12px;
  color: #9ca3af;
  font-size: 13px;
}

.btn-icon {
  font-size: 16px;
  margin-right: 4px;
}

.log-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.log-label {
  min-width: 100px;
  color: #6b7280;
  font-size: 14px;
}

.settings-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e9e4f5;
}

.rewrite-rules-container {
  padding: 20px 0;
}

.rewrite-toolbar {
  margin-bottom: 16px;
}

.rewrite-help {
  margin-top: 20px;
  padding: 16px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.rewrite-help ul {
  margin: 8px 0;
  padding-left: 20px;
}

.rewrite-help li {
  margin: 4px 0;
}
</style>
