<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-bold text-slate-900">Laporan</h2>
      <p class="text-sm text-slate-500">Ringkasan laporan manufaktur dan KPI.</p>
    </div>
    <div v-if="reportError" class="rounded-3xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ reportError }}</div>
    <div class="grid gap-6 lg:grid-cols-3">
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <h3 class="text-sm text-slate-500">Laporan Pendapatan</h3>
        <p class="mt-4 text-3xl font-bold text-slate-900">{{ overview ? formatCurrency(overview.totalRevenue) : 'Rp 0' }}</p>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <h3 class="text-sm text-slate-500">Laporan Produksi</h3>
        <p class="mt-4 text-3xl font-bold text-slate-900">{{ overview?.totalProduction ?? 0 }}</p>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <h3 class="text-sm text-slate-500">Laporan Pelanggan</h3>
        <p class="mt-4 text-3xl font-bold text-slate-900">{{ overview?.totalCustomers ?? 0 }}</p>
      </div>
    </div>
    <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
      <p class="text-sm text-slate-500">Bagian laporan ini akan menampilkan grafik analitik dan data ekspor ketika tersedia.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface AdminOverview {
  totalProduction: number
  activeOrders: number
  totalCustomers: number
  totalRevenue: number
  supportTickets: number
  pendingApprovals: number
}

const router = useRouter()
const overview = ref<AdminOverview | null>(null)
const reportError = ref('')

const fetchReportData = async () => {
  reportError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/admin/overview', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    overview.value = await response.json()
  } catch {
    reportError.value = 'Gagal memuat data laporan'
  }
}

const formatCurrency = (value: number) => `Rp ${value.toLocaleString('id-ID')}`

onMounted(() => {
  fetchReportData()
})
</script>

<style scoped>
</style>
