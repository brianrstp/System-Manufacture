import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '../pages/LandingPage.vue'
import LoginPage from '../pages/LoginPage.vue'
import AdminLayout from '../pages/admin/AdminLayout.vue'
import AdminDashboard from '../pages/admin/dashboard/AdminDashboard.vue'
import AdminProducts from '../pages/admin/products/ProductsPage.vue'
import AdminInventory from '../pages/admin/inventory/InventoryPage.vue'
import AdminOrders from '../pages/admin/orders/OrdersPage.vue'
import AdminCustomers from '../pages/admin/customers/CustomersPage.vue'
import AdminProduction from '../pages/admin/production/ProductionPage.vue'
import AdminCategories from '../pages/admin/categories/CategoriesPage.vue'
import AdminUnits from '../pages/admin/units/UnitsPage.vue'
import AdminWarehouses from '../pages/admin/warehouses/WarehousesPage.vue'
import AdminBOMs from '../pages/admin/boms/BOMsPage.vue'
import AdminStockMovements from '../pages/admin/stock-movements/StockMovementsPage.vue'
import AdminReports from '../pages/admin/reports/ReportsPage.vue'
import AdminSettings from '../pages/admin/settings/SettingsPage.vue'
import CustomerLoginPage from '../pages/customer/CustomerLoginPage.vue'
import CustomerOrdersPage from '../pages/customer/CustomerOrdersPage.vue'
import CustomerProfilePage from '../pages/customer/CustomerProfilePage.vue'
import CustomerHelpPage from '../pages/customer/CustomerHelpPage.vue'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: LandingPage,
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
  },
  {
    path: '/customer/login',
    name: 'CustomerLogin',
    component: CustomerLoginPage,
  },
  {
    path: '/customer/orders',
    name: 'CustomerOrders',
    component: CustomerOrdersPage,
    meta: { requiresCustomerAuth: true },
  },
  {
    path: '/customer/profile',
    name: 'CustomerProfile',
    component: CustomerProfilePage,
    meta: { requiresCustomerAuth: true },
  },
  {
    path: '/customer/help',
    name: 'CustomerHelp',
    component: CustomerHelpPage,
    meta: { requiresCustomerAuth: true },
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: AdminDashboard,
      },
      {
        path: 'products',
        name: 'AdminProducts',
        component: AdminProducts,
      },
      {
        path: 'inventory',
        name: 'AdminInventory',
        component: AdminInventory,
      },
      {
        path: 'boms',
        name: 'AdminBOMs',
        component: AdminBOMs,
      },
      {
        path: 'stock-movements',
        name: 'AdminStockMovements',
        component: AdminStockMovements,
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: AdminOrders,
      },
      {
        path: 'customers',
        name: 'AdminCustomers',
        component: AdminCustomers,
      },
      {
        path: 'production',
        name: 'AdminProduction',
        component: AdminProduction,
      },
      {
        path: 'reports',
        name: 'AdminReports',
        component: AdminReports,
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: AdminCategories,
      },
      {
        path: 'units',
        name: 'AdminUnits',
        component: AdminUnits,
      },
      {
        path: 'warehouses',
        name: 'AdminWarehouses',
        component: AdminWarehouses,
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: AdminSettings,
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const adminToken = localStorage.getItem('admin_token')
  const customerToken = localStorage.getItem('customer_token')

  if (to.meta.requiresAuth && !adminToken) {
    return { name: 'Login' }
  }

  if (to.meta.requiresCustomerAuth && !customerToken) {
    return { name: 'CustomerLogin' }
  }

  if (to.name === 'Login' && adminToken) {
    return { name: 'AdminDashboard' }
  }

  if (to.name === 'CustomerLogin' && customerToken) {
    return { name: 'CustomerOrders' }
  }
})

export default router
