<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Produksi</h2>
        <p class="text-sm text-slate-500">Pantau pekerjaan produksi dan performa manufaktur.</p>
      </div>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <input v-model="productionSearch" type="text" placeholder="Cari pekerjaan..." class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
        <select v-model="productionStatus" class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
          <option value="">Semua Status</option>
          <option value="pending">Pending</option>
          <option value="diproses">Diproses</option>
          <option value="selesai">Selesai</option>
          <option value="quality check">Quality Check</option>
          <option value="dikirim">Dikirim</option>
        </select>
        <button @click="fetchProductionJobs" class="rounded-full border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-50 transition">Segarkan</button>
        <button @click="openProductionForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Produksi</button>
      </div>
    </div>

    <div v-if="productionError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ productionError }}</div>

    <div class="grid gap-6 lg:grid-cols-3">
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Pekerjaan Produksi</p>
        <h3 class="mt-3 text-3xl font-bold text-slate-900">{{ productionJobs.length }}</h3>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Pekerjaan Selesai</p>
        <h3 class="mt-3 text-3xl font-bold text-slate-900">{{ completedJobs }}</h3>
      </div>
      <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
        <p class="text-sm text-slate-500">Pekerjaan Tertunda</p>
        <h3 class="mt-3 text-3xl font-bold text-slate-900">{{ pendingJobs }}</h3>
      </div>
    </div>

    <div class="rounded-3xl bg-white border border-slate-200 p-6 shadow-sm">
      <div class="flex items-center justify-between mb-5">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">Daftar Pekerjaan Produksi</h3>
          <p class="text-sm text-slate-500">Ringkasan pekerjaan manufaktur</p>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-left text-sm">
          <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
            <tr>
              <th class="px-6 py-4">Kode</th>
              <th class="px-6 py-4">Produk</th>
              <th class="px-6 py-4">Mulai</th>
              <th class="px-6 py-4">Durasi</th>
              <th class="px-6 py-4">Status</th>
              <th class="px-6 py-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr v-for="job in productionJobs" :key="job.id" class="hover:bg-slate-50 transition-colors">
              <td class="px-6 py-4">{{ job.jobCode }}</td>
              <td class="px-6 py-4">{{ job.product }}</td>
              <td class="px-6 py-4">{{ job.startDate }}</td>
              <td class="px-6 py-4">{{ job.durationDays }} hari</td>
              <td class="px-6 py-4">{{ job.status }}</td>
              <td class="px-6 py-4 text-right">
                <button @click="openProductionForm(job)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
                <button @click="deleteProductionJob(job.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showProductionForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeProductionForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingProduction ? 'Edit Pekerjaan Produksi' : 'Tambah Pekerjaan Produksi' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="productionFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ productionFormError }}</div>
          <div v-if="productionFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ productionFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Kode Pekerjaan</span>
            <input v-model="productionForm.jobCode" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productionFormErrors.jobCode" class="text-xs text-red-600">{{ productionFormErrors.jobCode }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Produk</span>
            <input v-model="productionForm.product" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productionFormErrors.product" class="text-xs text-red-600">{{ productionFormErrors.product }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Tanggal Mulai</span>
            <input v-model="productionForm.startDate" type="date" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productionFormErrors.startDate" class="text-xs text-red-600">{{ productionFormErrors.startDate }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Durasi (hari)</span>
            <input v-model.number="productionForm.durationDays" type="number" min="1" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productionFormErrors.durationDays" class="text-xs text-red-600">{{ productionFormErrors.durationDays }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Status</span>
            <select v-model="productionForm.status" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="pending">Pending</option>
              <option value="Diproses">Diproses</option>
              <option value="Selesai">Selesai</option>
              <option value="Quality Check">Quality Check</option>
              <option value="Dikirim">Dikirim</option>
            </select>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeProductionForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveProductionJob" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface ProductionJob {
  id: number
  jobCode: string
  product: string
  startDate: string
  durationDays: number
  status: string
}

const router = useRouter()
const productionJobs = ref<ProductionJob[]>([])
const productionError = ref('')
const productionSearch = ref('')
const productionStatus = ref('')
const showProductionForm = ref(false)
const editingProduction = ref(false)
const productionFormError = ref('')
const productionFormSuccess = ref('')
const productionFormErrors = ref<Record<string, string>>({})
const productionForm = ref({ id: 0, jobCode: '', product: '', startDate: new Date().toISOString().slice(0, 10), durationDays: 1, status: 'pending' })

const fetchProductionJobs = async () => {
  productionError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const params = new URLSearchParams()
    if (productionSearch.value.trim()) params.set('search', productionSearch.value.trim())
    if (productionStatus.value) params.set('status', productionStatus.value)
    const url = `/api/production${params.toString() ? `?${params.toString()}` : ''}`
    const response = await fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    productionJobs.value = (result.data || []).map((job: any) => ({
      ...job,
      startDate: typeof job.startDate === 'string' ? job.startDate.slice(0, 10) : job.startDate,
    }))
  } catch {
    productionError.value = 'Gagal memuat produksi'
  }
}

const openProductionForm = (job?: ProductionJob) => {
  productionFormError.value = ''
  productionFormSuccess.value = ''
  productionFormErrors.value = {}

  if (job) {
    editingProduction.value = true
    productionForm.value = {
      id: job.id,
      jobCode: job.jobCode,
      product: job.product,
      startDate: job.startDate,
      durationDays: job.durationDays,
      status: job.status,
    }
  } else {
    editingProduction.value = false
    productionForm.value = { id: 0, jobCode: '', product: '', startDate: new Date().toISOString().slice(0, 10), durationDays: 1, status: 'pending' }
  }
  showProductionForm.value = true
}

const closeProductionForm = () => {
  showProductionForm.value = false
  productionFormError.value = ''
  productionFormSuccess.value = ''
  productionFormErrors.value = {}
}

const validateProductionForm = () => {
  productionFormError.value = ''
  productionFormSuccess.value = ''
  productionFormErrors.value = {}

  if (!productionForm.value.jobCode.trim()) productionFormErrors.value.jobCode = 'Kode pekerjaan wajib diisi.'
  if (!productionForm.value.product.trim()) productionFormErrors.value.product = 'Nama produk wajib diisi.'
  if (!productionForm.value.startDate) productionFormErrors.value.startDate = 'Tanggal mulai wajib diisi.'
  if (!productionForm.value.durationDays || productionForm.value.durationDays <= 0) productionFormErrors.value.durationDays = 'Durasi harus lebih dari 0.'

  if (Object.keys(productionFormErrors.value).length > 0) {
    productionFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveProductionJob = async () => {
  if (!validateProductionForm()) return

  productionError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingProduction.value ? `/api/production/${productionForm.value.id}` : '/api/production'
  const method = editingProduction.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(productionForm.value) })
    if (!response.ok) throw new Error()
    await fetchProductionJobs()
    productionFormSuccess.value = 'Pekerjaan produksi berhasil disimpan.'
    setTimeout(() => closeProductionForm(), 800)
  } catch {
    productionError.value = 'Gagal menyimpan pekerjaan produksi'
  }
}

const deleteProductionJob = async (id: number) => {
  productionError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus pekerjaan produksi ini?')) return
  try {
    const response = await fetch(`/api/production/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchProductionJobs()
  } catch {
    productionError.value = 'Gagal menghapus pekerjaan produksi'
  }
}

const completedJobs = computed(() => productionJobs.value.filter((job) => job.status.toLowerCase() === 'selesai').length)
const pendingJobs = computed(() => productionJobs.value.filter((job) => job.status.toLowerCase() !== 'selesai').length)

onMounted(() => {
  fetchProductionJobs()
})
</script>

<style scoped>
</style>
