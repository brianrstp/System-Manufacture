<template>
  <div class="min-h-screen bg-slate-50 text-slate-900 font-sans">
    <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-black/40 z-40 lg:hidden"></div>

    <aside :class="['fixed top-0 left-0 h-full bg-white border-r border-slate-200 z-50 shadow-xl transition-transform duration-300 ease-in-out', sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0']" style="width: 280px;">
      <div class="h-full flex flex-col">
        <div class="p-6 border-b border-slate-200">
          <div class="text-primary-600 font-bold text-xl">Admin Portal</div>
          <p class="text-sm text-slate-500 mt-1">Kelola manufaktur dan operasi.</p>
        </div>

        <nav class="flex-1 p-6 overflow-y-auto">
          <ul class="space-y-2">
            <li v-for="item in navItems" :key="item.route">
              <RouterLink :to="{ name: item.route }" class="w-full block rounded-2xl px-4 py-3 transition hover:bg-slate-100" :class="item.route === currentRoute ? 'bg-slate-100 font-semibold text-slate-900' : 'text-slate-600'">
                <span class="inline-flex items-center gap-3">
                  <i :class="[item.icon, 'w-5 text-center']"></i>
                  {{ item.label }}
                </span>
              </RouterLink>
            </li>
          </ul>
        </nav>

        <div class="p-6 border-t border-slate-200">
          <div class="text-sm text-slate-500 mb-3">Masuk sebagai</div>
          <div class="flex items-center justify-between gap-3">
            <div>
              <div class="font-semibold text-slate-900">Admin Manufaktur</div>
              <div class="text-xs text-slate-500">admin@indomfg.co.id</div>
            </div>
            <button @click="logout" class="rounded-full bg-red-50 text-red-700 px-4 py-2 text-xs font-semibold hover:bg-red-100 transition">Logout</button>
          </div>
        </div>
      </div>
    </aside>

    <div class="lg:pl-[280px]">
      <header class="sticky top-0 z-30 bg-white border-b border-slate-200 backdrop-blur-xl bg-white/90">
        <div class="flex items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8">
          <div class="flex items-center gap-4">
            <button @click="sidebarOpen = !sidebarOpen" class="lg:hidden inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:text-slate-900 transition">
              <span class="sr-only">Toggle sidebar</span>
              <i class="fas fa-bars"></i>
            </button>
            <div>
              <p class="text-sm text-primary-600 font-semibold">Admin Portal</p>
              <h1 class="text-2xl font-bold text-slate-900">{{ currentPage }}</h1>
            </div>
          </div>
          <div class="hidden sm:flex items-center gap-3">
            <div class="text-right">
              <div class="text-sm font-semibold text-slate-900">Admin</div>
              <div class="text-xs text-slate-500">Panel Manufaktur</div>
            </div>
            <button @click="logout" class="rounded-2xl bg-red-50 px-4 py-2 text-sm font-semibold text-red-700 hover:bg-red-100 transition">Logout</button>
          </div>
        </div>
      </header>

      <main class="px-4 py-8 sm:px-6 lg:px-10">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute, RouterLink, RouterView } from 'vue-router'

const router = useRouter()
const route = useRoute()
const sidebarOpen = ref(false)

const navItems = [
  { route: 'AdminDashboard', label: 'Dashboard', icon: 'fas fa-home' },
  { route: 'AdminProducts', label: 'Produk', icon: 'fas fa-box-open' },
  { route: 'AdminInventory', label: 'Inventaris', icon: 'fas fa-warehouse' },
  { route: 'AdminBOMs', label: 'BOM', icon: 'fas fa-layer-group' },
  { route: 'AdminStockMovements', label: 'Mutasi Stok', icon: 'fas fa-exchange-alt' },
  { route: 'AdminOrders', label: 'Pesanan', icon: 'fas fa-clipboard-list' },
  { route: 'AdminCustomers', label: 'Pelanggan', icon: 'fas fa-user-friends' },
  { route: 'AdminProduction', label: 'Produksi', icon: 'fas fa-industry' },
  { route: 'AdminReports', label: 'Laporan', icon: 'fas fa-chart-line' },
  { route: 'AdminCategories', label: 'Kategori', icon: 'fas fa-tags' },
  { route: 'AdminUnits', label: 'Unit', icon: 'fas fa-balance-scale' },
  { route: 'AdminWarehouses', label: 'Gudang', icon: 'fas fa-warehouse' },
  { route: 'AdminSettings', label: 'Pengaturan', icon: 'fas fa-cog' },
]

const currentRoute = computed(() => route.name as string)
const currentPage = computed(() => navItems.find((item) => item.route === currentRoute.value)?.label || 'Admin')

const logout = () => {
  localStorage.removeItem('admin_token')
  router.push({ name: 'Login' })
}
</script>

<style scoped>
.sidebar-transition { transition: transform 0.3s ease, width 0.3s ease; }
</style>
