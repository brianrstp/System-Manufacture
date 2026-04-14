<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Kategori</h2>
        <p class="text-sm text-slate-500">Kelola kategori produk dan struktur grup.</p>
      </div>
      <button @click="openCategoryForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Kategori</button>
    </div>

    <div v-if="categoryError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ categoryError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Nama</th>
            <th class="px-6 py-4">Slug</th>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="category in categories" :key="category.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ category.name }}</td>
            <td class="px-6 py-4">{{ category.slug }}</td>
            <td class="px-6 py-4">{{ category.status }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openCategoryForm(category)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
              <button @click="deleteCategory(category.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showCategoryForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeCategoryForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingCategory ? 'Edit Kategori' : 'Tambah Kategori' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="categoryFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ categoryFormError }}</div>
          <div v-if="categoryFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ categoryFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Nama</span>
            <input v-model="categoryForm.name" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="categoryFormErrors.name" class="text-xs text-red-600">{{ categoryFormErrors.name }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Slug</span>
            <input v-model="categoryForm.slug" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="categoryFormErrors.slug" class="text-xs text-red-600">{{ categoryFormErrors.slug }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Deskripsi</span>
            <textarea v-model="categoryForm.description" rows="4" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none"></textarea>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Status</span>
            <select v-model="categoryForm.status" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeCategoryForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveCategory" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface Category {
  id: number
  name: string
  slug: string
  description: string
  status: string
}

const router = useRouter()
const categories = ref<Category[]>([])
const categoryError = ref('')
const categoryFormError = ref('')
const categoryFormSuccess = ref('')
const categoryFormErrors = ref<Record<string, string>>({})
const showCategoryForm = ref(false)
const editingCategory = ref(false)
const categoryForm = ref({ id: 0, name: '', slug: '', description: '', status: 'active' })

const fetchCategories = async () => {
  categoryError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/categories', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    categories.value = result.data || []
  } catch {
    categoryError.value = 'Gagal memuat kategori'
  }
}

const openCategoryForm = (category?: Category) => {
  categoryFormError.value = ''
  categoryFormSuccess.value = ''
  categoryFormErrors.value = {}

  if (category) {
    editingCategory.value = true
    categoryForm.value = { ...category }
  } else {
    editingCategory.value = false
    categoryForm.value = { id: 0, name: '', slug: '', description: '', status: 'active' }
  }
  showCategoryForm.value = true
}

const closeCategoryForm = () => {
  showCategoryForm.value = false
  categoryFormError.value = ''
  categoryFormSuccess.value = ''
  categoryFormErrors.value = {}
}

const validateCategoryForm = () => {
  categoryFormError.value = ''
  categoryFormSuccess.value = ''
  categoryFormErrors.value = {}

  if (!categoryForm.value.name.trim()) categoryFormErrors.value.name = 'Nama kategori wajib diisi.'
  if (!categoryForm.value.slug.trim()) categoryFormErrors.value.slug = 'Slug kategori wajib diisi.'

  if (Object.keys(categoryFormErrors.value).length > 0) {
    categoryFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveCategory = async () => {
  if (!validateCategoryForm()) return

  categoryError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingCategory.value ? `/api/categories/${categoryForm.value.id}` : '/api/categories'
  const method = editingCategory.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(categoryForm.value) })
    if (!response.ok) throw new Error()
    await fetchCategories()
    categoryFormSuccess.value = 'Kategori berhasil disimpan.'
    setTimeout(() => closeCategoryForm(), 800)
  } catch {
    categoryError.value = 'Gagal menyimpan kategori'
  }
}

const deleteCategory = async (id: number) => {
  categoryError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus kategori ini?')) return
  try {
    const response = await fetch(`/api/categories/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchCategories()
  } catch {
    categoryError.value = 'Gagal menghapus kategori'
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
</style>
