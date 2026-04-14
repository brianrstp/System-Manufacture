<template>
  <div class="space-y-8">
    <section class="grid gap-6 lg:grid-cols-4">
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Total Pesanan</p>
        <h2 class="mt-3 text-3xl font-bold text-slate-900">{{ overview?.activeOrders ?? 0 }}</h2>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Pendapatan</p>
        <h2 class="mt-3 text-3xl font-bold text-slate-900">Rp {{ formatNumber(overview?.totalRevenue ?? 0) }}</h2>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Produksi Aktif</p>
        <h2 class="mt-3 text-3xl font-bold text-slate-900">{{ overview?.totalProduction ?? 0 }}</h2>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Pelanggan</p>
        <h2 class="mt-3 text-3xl font-bold text-slate-900">{{ overview?.totalCustomers ?? 0 }}</h2>
      </div>
    </section>

    <section class="grid gap-6 lg:grid-cols-2">
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <div class="flex items-center justify-between mb-5">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">Pendapatan Bulanan</h3>
            <p class="text-sm text-slate-500">Grafik pendapatan terbaru</p>
          </div>
        </div>
        <div class="relative h-72">
          <canvas id="salesChart"></canvas>
        </div>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <div class="flex items-center justify-between mb-5">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">Status Produksi</h3>
            <p class="text-sm text-slate-500">Ringkasan proses saat ini</p>
          </div>
        </div>
        <div class="relative h-72">
          <canvas id="statusChart"></canvas>
        </div>
      </div>
    </section>

    <section class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-xl font-bold text-slate-900">Pesanan Terbaru</h3>
          <p class="text-sm text-slate-500">Top order terbaru dari pelanggan.</p>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-left text-sm border-separate border-spacing-0">
          <thead class="bg-slate-50 text-slate-500 text-xs uppercase tracking-[0.2em]">
            <tr>
              <th class="px-5 py-4">ID</th>
              <th class="px-5 py-4">Pelanggan</th>
              <th class="px-5 py-4">Produk</th>
              <th class="px-5 py-4">Tanggal</th>
              <th class="px-5 py-4">Jumlah</th>
              <th class="px-5 py-4">Status</th>
            </tr>
          </thead>
          <tbody class="text-slate-700">
            <tr v-for="order in recentOrders" :key="order.id" class="border-t border-slate-100 hover:bg-slate-50">
              <td class="px-5 py-4 font-mono text-slate-900">#{{ order.id }}</td>
              <td class="px-5 py-4">{{ order.customerName }}</td>
              <td class="px-5 py-4">{{ order.product }}</td>
              <td class="px-5 py-4">{{ order.orderDate }}</td>
              <td class="px-5 py-4">{{ order.amount }}</td>
              <td class="px-5 py-4">
                <span :class="['inline-flex rounded-full px-3 py-1 text-xs font-semibold', statusClasses(order.status)]">{{ order.status }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import Chart from 'chart.js/auto'
import { useRouter } from 'vue-router'

interface Order {
  id: number
  customerName: string
  product: string
  orderDate: string
  amount: string
  status: string
}

interface AdminOverview {
  totalProduction: number
  activeOrders: number
  totalCustomers: number
  totalRevenue: number
  monthlyRevenue: number[]
  productionStatusCounts: Record<string, number>
  supportTickets: number
  pendingApprovals: number
}

const router = useRouter()
const overview = ref<AdminOverview | null>(null)
const orders = ref<Order[]>([])
let salesChart: Chart | null = null
let statusChart: Chart | null = null

const fetchAdminOverview = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })

  try {
    const response = await fetch('/api/admin/overview', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error('Network response was not ok')
    overview.value = await response.json()
  } catch {
    overview.value = {
      totalProduction: 0,
      activeOrders: 0,
      totalCustomers: 0,
      totalRevenue: 0,
      monthlyRevenue: [0, 0, 0, 0, 0, 0],
      productionStatusCounts: {},
      supportTickets: 0,
      pendingApprovals: 0,
    }
  }
}

const fetchRecentOrders = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })

  try {
    const response = await fetch('/api/orders?limit=5', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error('Network response was not ok')
    const result = await response.json()
    orders.value = result.data || []
  } catch {
    orders.value = []
  }
}

const recentOrders = computed(() => orders.value.slice(0, 5))

const statusClasses = (status: string) => {
  switch (status) {
    case 'Diproses': return 'bg-blue-100 text-blue-700'
    case 'Selesai': return 'bg-green-100 text-green-700'
    case 'Quality Check': return 'bg-amber-100 text-amber-700'
    case 'Dikirim': return 'bg-purple-100 text-purple-700'
    case 'Pending': return 'bg-slate-100 text-slate-700'
    default: return 'bg-slate-100 text-slate-700'
  }
}

const formatNumber = (value: number | string) => {
  const num = typeof value === 'string' ? Number(value) : value
  return num.toLocaleString('id-ID')
}

const initCharts = () => {
  const salesCanvas = document.getElementById('salesChart') as HTMLCanvasElement | null
  const statusCanvas = document.getElementById('statusChart') as HTMLCanvasElement | null

  const revenueData = overview.value?.monthlyRevenue ?? [0, 0, 0, 0, 0, 0]
  const revenueLabels = ['-5 bln', '-4 bln', '-3 bln', '-2 bln', '-1 bln', 'Saat ini']
  const statusCounts = overview.value?.productionStatusCounts ?? {}
  const doughnutLabels = Object.keys(statusCounts).length
    ? Object.keys(statusCounts)
    : ['Diproses', 'Selesai', 'Quality Check', 'Pending']
  const doughnutData = Object.keys(statusCounts).length
    ? Object.values(statusCounts)
    : [1, 1, 1, 1]

  if (salesCanvas) {
    if (salesChart) salesChart.destroy()
    salesChart = new Chart(salesCanvas, {
      type: 'line',
      data: {
        labels: revenueLabels,
        datasets: [{ label: 'Pendapatan', data: revenueData, borderColor: '#2563EB', backgroundColor: 'rgba(37,99,235,0.12)', fill: true, tension: 0.4, pointRadius: 3 }],
      },
      options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { display: false }, ticks: { color: '#64748b' } }, y: { grid: { color: '#e2e8f0' }, ticks: { color: '#64748b', callback: (value) => `Rp ${Number(value) / 1000000}M` } } } },
    })
  }

  if (statusCanvas) {
    if (statusChart) statusChart.destroy()
    statusChart = new Chart(statusCanvas, {
      type: 'doughnut',
      data: {
        labels: doughnutLabels,
        datasets: [{ data: doughnutData, backgroundColor: ['#2563EB', '#16A34A', '#F59E0B', '#94A3B8', '#8B5CF6', '#EC4899'], borderWidth: 0 }],
      },
      options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { color: '#334155' } } } },
    })
  }
}

onMounted(async () => {
  await fetchAdminOverview()
  await fetchRecentOrders()
  initCharts()
})
</script>

<style scoped>
</style>
