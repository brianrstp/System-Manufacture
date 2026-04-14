<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Pelanggan</h2>
        <p class="text-sm text-slate-500">Kelola data pelanggan, kontak, dan status akun.</p>
      </div>
      <button @click="openCustomerForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Pelanggan</button>
    </div>

    <div v-if="customerError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ customerError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Nama</th>
            <th class="px-6 py-4">Email</th>
            <th class="px-6 py-4">Telepon</th>
            <th class="px-6 py-4">Alamat</th>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="customer in customers" :key="customer.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ customer.name }}</td>
            <td class="px-6 py-4">{{ customer.email }}</td>
            <td class="px-6 py-4">{{ customer.phone }}</td>
            <td class="px-6 py-4">{{ customer.address }}</td>
            <td class="px-6 py-4">{{ customer.status }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openCustomerForm(customer)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
              <button @click="deleteCustomer(customer.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showCustomerForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeCustomerForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingCustomer ? 'Edit Pelanggan' : 'Tambah Pelanggan' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="customerFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ customerFormError }}</div>
          <div v-if="customerFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ customerFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Nama</span>
            <input v-model="customerForm.name" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="customerFormErrors.name" class="text-xs text-red-600">{{ customerFormErrors.name }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Email</span>
            <input v-model="customerForm.email" type="email" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="customerFormErrors.email" class="text-xs text-red-600">{{ customerFormErrors.email }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Telepon</span>
            <input v-model="customerForm.phone" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Alamat</span>
            <textarea v-model="customerForm.address" rows="4" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none"></textarea>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Kata Sandi</span>
            <input v-model="customerForm.password" type="password" placeholder="Kosongkan jika tidak ingin mengganti" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="customerFormErrors.password" class="text-xs text-red-600">{{ customerFormErrors.password }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Status</span>
            <select v-model="customerForm.status" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeCustomerForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveCustomer" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface Customer {
  id: number
  name: string
  email: string
  phone: string
  address: string
  status: string
}

const router = useRouter()
const customers = ref<Customer[]>([])
const customerError = ref('')
const customerFormError = ref('')
const customerFormSuccess = ref('')
const customerFormErrors = ref<Record<string, string>>({})
const showCustomerForm = ref(false)
const editingCustomer = ref(false)
const customerForm = ref({ id: 0, name: '', email: '', phone: '', address: '', password: '', status: 'active' })

const fetchCustomers = async () => {
  customerError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/customers', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    customers.value = result.data || []
  } catch {
    customerError.value = 'Gagal memuat pelanggan'
  }
}

const openCustomerForm = (customer?: Customer) => {
  customerFormError.value = ''
  customerFormSuccess.value = ''
  customerFormErrors.value = {}

  if (customer) {
    editingCustomer.value = true
    customerForm.value = { ...customer, password: '' }
  } else {
    editingCustomer.value = false
    customerForm.value = { id: 0, name: '', email: '', phone: '', address: '', password: '', status: 'active' }
  }
  showCustomerForm.value = true
}

const closeCustomerForm = () => {
  showCustomerForm.value = false
  customerFormError.value = ''
  customerFormSuccess.value = ''
  customerFormErrors.value = {}
}

const validateCustomerForm = () => {
  customerFormError.value = ''
  customerFormSuccess.value = ''
  customerFormErrors.value = {}
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

  if (!customerForm.value.name.trim()) customerFormErrors.value.name = 'Nama wajib diisi.'
  if (!customerForm.value.email.trim()) customerFormErrors.value.email = 'Email wajib diisi.'
  else if (!emailRegex.test(customerForm.value.email.trim())) customerFormErrors.value.email = 'Email tidak valid.'
  if (customerForm.value.password.trim() && customerForm.value.password.trim().length < 8) customerFormErrors.value.password = 'Kata sandi minimal 8 karakter jika diisi.'

  if (Object.keys(customerFormErrors.value).length > 0) {
    customerFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveCustomer = async () => {
  if (!validateCustomerForm()) return

  customerError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingCustomer.value ? `/api/customers/${customerForm.value.id}` : '/api/customers'
  const method = editingCustomer.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(customerForm.value) })
    if (!response.ok) throw new Error()
    await fetchCustomers()
    customerFormSuccess.value = 'Pelanggan berhasil disimpan.'
    setTimeout(() => closeCustomerForm(), 800)
  } catch {
    customerError.value = 'Gagal menyimpan pelanggan'
  }
}

const deleteCustomer = async (id: number) => {
  customerError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus pelanggan ini?')) return
  try {
    const response = await fetch(`/api/customers/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchCustomers()
  } catch {
    customerError.value = 'Gagal menghapus pelanggan'
  }
}

onMounted(() => {
  fetchCustomers()
})
</script>

<style scoped>
</style>
