<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Gudang</h2>
        <p class="text-sm text-slate-500">Kelola lokasi penyimpanan dan informasi gudang.</p>
      </div>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <input v-model="warehouseSearch" type="text" placeholder="Cari gudang..." class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
        <select v-model="warehouseStatus" class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
          <option value="">Semua Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
        <button @click="fetchWarehouses" class="rounded-full bg-slate-100 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-200 transition">Refresh</button>
        <button @click="openWarehouseForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Gudang</button>
      </div>
    </div>

    <div v-if="warehouseError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ warehouseError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Kode</th>
            <th class="px-6 py-4">Nama</th>
            <th class="px-6 py-4">Lokasi</th>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="warehouse in warehouses" :key="warehouse.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ warehouse.code }}</td>
            <td class="px-6 py-4">{{ warehouse.name }}</td>
            <td class="px-6 py-4">{{ warehouse.location }}</td>
            <td class="px-6 py-4">{{ warehouse.status }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openWarehouseForm(warehouse)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
              <button @click="deleteWarehouse(warehouse.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showWarehouseForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeWarehouseForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingWarehouse ? 'Edit Gudang' : 'Tambah Gudang' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="warehouseFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ warehouseFormError }}</div>
          <div v-if="warehouseFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ warehouseFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Kode</span>
            <input v-model="warehouseForm.code" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="warehouseFormErrors.code" class="text-xs text-red-600">{{ warehouseFormErrors.code }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Nama</span>
            <input v-model="warehouseForm.name" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="warehouseFormErrors.name" class="text-xs text-red-600">{{ warehouseFormErrors.name }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Lokasi</span>
            <input v-model="warehouseForm.location" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="warehouseFormErrors.location" class="text-xs text-red-600">{{ warehouseFormErrors.location }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Deskripsi</span>
            <textarea v-model="warehouseForm.description" rows="4" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none"></textarea>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Status</span>
            <select v-model="warehouseForm.status" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeWarehouseForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveWarehouse" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface Warehouse {
  id: number
  code: string
  name: string
  description: string
  location: string
  status: string
}

const router = useRouter()
const warehouses = ref<Warehouse[]>([])
const warehouseSearch = ref('')
const warehouseStatus = ref('')
const warehouseError = ref('')
const warehouseFormError = ref('')
const warehouseFormSuccess = ref('')
const warehouseFormErrors = ref<Record<string, string>>({})
const showWarehouseForm = ref(false)
const editingWarehouse = ref(false)
const warehouseForm = ref({ id: 0, code: '', name: '', description: '', location: '', status: 'active' })

const fetchWarehouses = async () => {
  warehouseError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const params = new URLSearchParams()
    if (warehouseSearch.value.trim()) params.set('search', warehouseSearch.value.trim())
    if (warehouseStatus.value) params.set('status', warehouseStatus.value)
    const url = `/api/warehouses${params.toString() ? `?${params.toString()}` : ''}`
    const response = await fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    warehouses.value = result.data || []
  } catch {
    warehouseError.value = 'Gagal memuat gudang'
  }
}

const openWarehouseForm = (warehouse?: Warehouse) => {
  warehouseFormError.value = ''
  warehouseFormSuccess.value = ''
  warehouseFormErrors.value = {}

  if (warehouse) {
    editingWarehouse.value = true
    warehouseForm.value = { ...warehouse }
  } else {
    editingWarehouse.value = false
    warehouseForm.value = { id: 0, code: '', name: '', description: '', location: '', status: 'active' }
  }
  showWarehouseForm.value = true
}

const closeWarehouseForm = () => {
  showWarehouseForm.value = false
  warehouseFormError.value = ''
  warehouseFormSuccess.value = ''
  warehouseFormErrors.value = {}
}

const validateWarehouseForm = () => {
  warehouseFormError.value = ''
  warehouseFormSuccess.value = ''
  warehouseFormErrors.value = {}

  if (!warehouseForm.value.code.trim()) warehouseFormErrors.value.code = 'Kode gudang wajib diisi.'
  if (!warehouseForm.value.name.trim()) warehouseFormErrors.value.name = 'Nama gudang wajib diisi.'
  if (!warehouseForm.value.location.trim()) warehouseFormErrors.value.location = 'Lokasi gudang wajib diisi.'

  if (Object.keys(warehouseFormErrors.value).length > 0) {
    warehouseFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveWarehouse = async () => {
  if (!validateWarehouseForm()) return

  warehouseError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingWarehouse.value ? `/api/warehouses/${warehouseForm.value.id}` : '/api/warehouses'
  const method = editingWarehouse.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(warehouseForm.value) })
    if (!response.ok) throw new Error()
    await fetchWarehouses()
    warehouseFormSuccess.value = 'Gudang berhasil disimpan.'
    setTimeout(() => closeWarehouseForm(), 800)
  } catch {
    warehouseError.value = 'Gagal menyimpan gudang'
  }
}

const deleteWarehouse = async (id: number) => {
  warehouseError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus gudang ini?')) return
  try {
    const response = await fetch(`/api/warehouses/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchWarehouses()
  } catch {
    warehouseError.value = 'Gagal menghapus gudang'
  }
}

onMounted(() => {
  fetchWarehouses()
})
</script>

<style scoped>
</style>
