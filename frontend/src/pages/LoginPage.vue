<template>
  <main class="min-h-screen bg-slate-950 text-white flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="w-full max-w-md space-y-8 rounded-3xl bg-slate-900/90 p-10 shadow-2xl shadow-slate-950/40 backdrop-blur">
      <div class="text-center">
        <p class="text-sm font-semibold uppercase tracking-[0.3em] text-primary-300">Admin Login</p>
        <h1 class="mt-4 text-3xl font-extrabold">Masuk ke Dashboard Admin</h1>
        <p class="mt-3 text-sm text-slate-400">Gunakan akun admin Anda untuk mengelola data dan laporan.</p>
      </div>

      <form @submit.prevent="login" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-slate-300">Username</label>
          <input v-model="username" type="text" required autocomplete="username"
                 class="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-950/90 px-4 py-3 text-white outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20" />
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-300">Password</label>
          <input v-model="password" type="password" required autocomplete="current-password"
                 class="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-950/90 px-4 py-3 text-white outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20" />
        </div>

        <div v-if="errorMessage" class="rounded-2xl bg-red-500/10 border border-red-500/20 px-4 py-3 text-sm text-red-200">
          {{ errorMessage }}
        </div>

        <button type="submit" :disabled="loading"
                class="w-full rounded-2xl bg-primary-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-primary-500 disabled:cursor-not-allowed disabled:opacity-60">
          <span v-if="!loading">Masuk</span>
          <span v-else>Memproses...</span>
        </button>
      </form>

      <p class="text-center text-sm text-slate-500">
        Gunakan <span class="font-semibold text-primary-300">admin</span> / <span class="font-semibold text-primary-300">admin123</span> jika belum ada akun.
      </p>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')
const router = useRouter()

const login = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch('/api/admin/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username: username.value, password: password.value }),
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: 'Login gagal' }))
      errorMessage.value = error.message || 'Login gagal'
      return
    }

    const data = await response.json()
    localStorage.setItem('admin_token', data.token)
    router.push({ name: 'AdminDashboard' })
  } catch (err) {
    errorMessage.value = 'Tidak dapat terhubung ke server'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped></style>
