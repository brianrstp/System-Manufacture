<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Unit</h2>
        <p class="text-sm text-slate-500">Kelola unit konversi dan satuan produk.</p>
      </div>
      <button @click="openUnitForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Unit</button>
    </div>

    <div v-if="unitError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ unitError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Kode</th>
            <th class="px-6 py-4">Nama</th>
            <th class="px-6 py-4">Faktor</th>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="unit in units" :key="unit.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ unit.code }}</td>
            <td class="px-6 py-4">{{ unit.name }}</td>
            <td class="px-6 py-4">{{ unit.factor }}</td>
            <td class="px-6 py-4">{{ unit.status }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openUnitForm(unit)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
              <button @click="deleteUnit(unit.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showUnitForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeUnitForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingUnit ? 'Edit Unit' : 'Tambah Unit' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="unitFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ unitFormError }}</div>
          <div v-if="unitFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ unitFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Kode</span>
            <input v-model="unitForm.code" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="unitFormErrors.code" class="text-xs text-red-600">{{ unitFormErrors.code }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Nama</span>
            <input v-model="unitForm.name" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="unitFormErrors.name" class="text-xs text-red-600">{{ unitFormErrors.name }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Faktor</span>
            <input v-model.number="unitForm.factor" type="number" step="0.1" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="unitFormErrors.factor" class="text-xs text-red-600">{{ unitFormErrors.factor }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Status</span>
            <select v-model="unitForm.status" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeUnitForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveUnit" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface Unit {
  id: number
  code: string
  name: string
  description: string
  factor: number
  status: string
}

const router = useRouter()
const units = ref<Unit[]>([])
const unitError = ref('')
const unitFormError = ref('')
const unitFormSuccess = ref('')
const unitFormErrors = ref<Record<string, string>>({})
const showUnitForm = ref(false)
const editingUnit = ref(false)
const unitForm = ref({ id: 0, code: '', name: '', description: '', factor: 1.0, status: 'active' })

const fetchUnits = async () => {
  unitError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/units', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    units.value = result.data || []
  } catch {
    unitError.value = 'Gagal memuat unit'
  }
}

const openUnitForm = (unit?: Unit) => {
  unitFormError.value = ''
  unitFormSuccess.value = ''
  unitFormErrors.value = {}

  if (unit) {
    editingUnit.value = true
    unitForm.value = { ...unit }
  } else {
    editingUnit.value = false
    unitForm.value = { id: 0, code: '', name: '', description: '', factor: 1.0, status: 'active' }
  }
  showUnitForm.value = true
}

const closeUnitForm = () => {
  showUnitForm.value = false
  unitFormError.value = ''
  unitFormSuccess.value = ''
  unitFormErrors.value = {}
}

const validateUnitForm = () => {
  unitFormError.value = ''
  unitFormSuccess.value = ''
  unitFormErrors.value = {}

  if (!unitForm.value.code.trim()) unitFormErrors.value.code = 'Kode unit wajib diisi.'
  if (!unitForm.value.name.trim()) unitFormErrors.value.name = 'Nama unit wajib diisi.'
  if (unitForm.value.factor <= 0) unitFormErrors.value.factor = 'Faktor harus lebih besar dari nol.'

  if (Object.keys(unitFormErrors.value).length > 0) {
    unitFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveUnit = async () => {
  if (!validateUnitForm()) return

  unitError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingUnit.value ? `/api/units/${unitForm.value.id}` : '/api/units'
  const method = editingUnit.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(unitForm.value) })
    if (!response.ok) throw new Error()
    await fetchUnits()
    unitFormSuccess.value = 'Unit berhasil disimpan.'
    setTimeout(() => closeUnitForm(), 800)
  } catch {
    unitError.value = 'Gagal menyimpan unit'
  }
}

const deleteUnit = async (id: number) => {
  unitError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus unit ini?')) return
  try {
    const response = await fetch(`/api/units/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchUnits()
  } catch {
    unitError.value = 'Gagal menghapus unit'
  }
}

onMounted(() => {
  fetchUnits()
})
</script>

<style scoped>
</style>
