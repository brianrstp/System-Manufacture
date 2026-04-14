<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-black/40 z-40 lg:hidden"></div>

    <aside :class="['fixed top-0 left-0 h-full bg-white border-r border-slate-200 z-50 shadow-xl transition-transform duration-300 ease-in-out', sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0']" style="width: 280px;">
      <div class="h-full flex flex-col">
        <div class="p-6 border-b border-slate-200">
          <div class="text-primary-600 font-bold text-xl">Customer Portal</div>
          <p class="text-sm text-slate-500 mt-1">Akses pesanan dan akun Anda.</p>
        </div>
        <nav class="flex-1 p-6 overflow-y-auto">
          <ul class="space-y-2">
            <li v-for="item in navItems" :key="item.route">
              <button @click="navigate(item.route)" class="w-full text-left rounded-2xl px-4 py-3 transition hover:bg-slate-100" :class="item.route === router.currentRoute.value.name ? 'bg-slate-100 font-semibold text-slate-900' : 'text-slate-600'">
                {{ item.label }}
              </button>
            </li>
          </ul>
        </nav>
        <div class="p-6 border-t border-slate-200">
          <div class="text-sm text-slate-500 mb-3">Masuk sebagai</div>
          <div class="flex items-center justify-between gap-3">
            <div>
              <div class="font-semibold text-slate-900">{{ customerName }}</div>
              <div class="text-xs text-slate-500">Customer</div>
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
              <p class="text-sm text-primary-600 font-semibold">Customer Portal</p>
              <h1 class="text-2xl font-bold text-slate-900">Pesanan Saya</h1>
            </div>
          </div>
          <div class="hidden sm:flex items-center gap-3">
            <div class="text-right">
              <div class="text-sm font-semibold text-slate-900">{{ customerName }}</div>
              <div class="text-xs text-slate-500">Akun pelanggan</div>
            </div>
            <button @click="logout" class="rounded-2xl bg-red-50 px-4 py-2 text-sm font-semibold text-red-700 hover:bg-red-100 transition">Logout</button>
          </div>
        </div>
      </header>

      <main class="px-4 py-8 sm:px-6 lg:px-10">
        <section class="grid gap-4 lg:grid-cols-3 mb-8">
          <div class="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <p class="text-sm text-slate-500">Total Pesanan</p>
            <p class="mt-4 text-3xl font-bold text-slate-900">{{ orders.length }}</p>
          </div>
          <div class="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <p class="text-sm text-slate-500">Pesanan Selesai</p>
            <p class="mt-4 text-3xl font-bold text-slate-900">{{ statusCount('Selesai') }}</p>
          </div>
          <div class="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <p class="text-sm text-slate-500">Pesanan Dalam Proses</p>
            <p class="mt-4 text-3xl font-bold text-slate-900">{{ statusCount('Diproses') }}</p>
          </div>
        </section>

        <section class="rounded-3xl bg-white border border-slate-200 shadow-sm p-6">
          <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between mb-6">
            <div>
              <p class="text-sm font-medium text-slate-500">Ringkasan Pesanan</p>
              <h2 class="text-xl font-bold text-slate-900">Pesanan Saya</h2>
            </div>
            <div class="text-sm text-slate-500">{{ orders.length }} item</div>
          </div>

          <div v-if="errorMessage" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-6">
            {{ errorMessage }}
          </div>

          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-slate-200 text-sm">
              <thead class="bg-slate-50 text-slate-500 uppercase tracking-[0.15em] text-xs">
                <tr>
                  <th class="px-6 py-4 text-left font-semibold">No. Pesanan</th>
                  <th class="px-6 py-4 text-left font-semibold">Produk</th>
                  <th class="px-6 py-4 text-left font-semibold">Tanggal</th>
                  <th class="px-6 py-4 text-right font-semibold">Jumlah</th>
                  <th class="px-6 py-4 text-left font-semibold">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 bg-white">
                <tr v-for="order in orders" :key="order.id" class="hover:bg-slate-50 transition-colors">
                  <td class="px-6 py-4 font-medium text-slate-900">{{ order.orderNumber }}</td>
                  <td class="px-6 py-4 text-slate-600">{{ order.product }}</td>
                  <td class="px-6 py-4 text-slate-500">{{ formatDate(order.orderDate) }}</td>
                  <td class="px-6 py-4 text-right text-slate-900 font-semibold">{{ formatAmount(order.amount) }}</td>
                  <td class="px-6 py-4">
                    <span :class="['inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold', statusClasses(order.status)]">
                      {{ order.status }}
                    </span>
                  </td>
                </tr>
                <tr v-if="orders.length === 0">
                  <td colspan="5" class="px-6 py-8 text-center text-sm text-slate-500">Tidak ada pesanan untuk ditampilkan.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

type CustomerOrder = {
  id: number
  orderNumber: string
  customerId: number
  customerName: string
  product: string
  orderDate: string
  amount: number
  status: string
}

const orders = ref<CustomerOrder[]>([])
const errorMessage = ref('')
const sidebarOpen = ref(false)
const router = useRouter()
const customerName = localStorage.getItem('customer_name') || 'Customer'

const navItems = [
  { route: 'CustomerOrders', label: 'Pesanan Saya' },
  { route: 'CustomerProfile', label: 'Profil' },
  { route: 'CustomerHelp', label: 'Bantuan' },
]

const navigate = (routeName: string) => {
  router.push({ name: routeName })
}

const statusClasses = (status: string) => {
  const map: Record<string, string> = {
    pending: 'bg-gray-100 text-gray-700',
    Diproses: 'bg-blue-100 text-blue-700',
    Selesai: 'bg-green-100 text-green-700',
    'Quality Check': 'bg-yellow-100 text-yellow-700',
    Dikirim: 'bg-purple-100 text-purple-700',
  }
  return map[status] || 'bg-gray-100 text-gray-700'
}

const formatDate = (value: string) => {
  try {
    return new Date(value).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch {
    return value
  }
}

const formatAmount = (value: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
}

const logout = () => {
  localStorage.removeItem('customer_token')
  localStorage.removeItem('customer_id')
  localStorage.removeItem('customer_name')
  router.push({ name: 'CustomerLogin' })
}

const statusCount = (status: string) => orders.value.filter((order) => order.status === status).length

const fetchOrders = async () => {
  errorMessage.value = ''
  const token = localStorage.getItem('customer_token')
  if (!token) return router.push({ name: 'CustomerLogin' })

  try {
    const response = await fetch('/api/customer/orders', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) {
      if (response.status === 401) {
        logout()
        return
      }
      throw new Error('Network response was not ok')
    }
    const result = await response.json()
    orders.value = result.data || []
  } catch (error) {
    errorMessage.value = 'Gagal memuat pesanan Anda'
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped></style>
