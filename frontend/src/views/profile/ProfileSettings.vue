<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../../components/BaseButton.vue';
import BaseInput from '../../components/BaseInput.vue';
import BaseAlert from '../../components/BaseAlert.vue';
import BaseToggle from '../../components/BaseToggle.vue';
import BaseSelect from '../../components/BaseSelect.vue';
import BaseDatePicker from '../../components/BaseDatePicker.vue';

const authStore = useAuthStore();
const isSaving = ref(false);
const isEditing = ref(false);
const saveSuccess = ref(false);
const errorMessage = ref('');

// Form Data State
const formData = reactive({
  full_name: '',
  username: '',
  phone_number: '',
  gender: '',
  date_of_birth: '',
  bio: '',
  newsletter_subscribed: false
});

const genderOptions = [
  { value: 'male', label: 'Male' },
  { value: 'female', label: 'Female' },
  { value: 'other', label: 'Prefer not to say' }
];

// Avatar state
const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const previewUrl = ref<string | null>(null);

const populateForm = () => {
  if (authStore.user) {
    formData.full_name = authStore.user.full_name || '';
    formData.username = authStore.user.username || '';
    formData.phone_number = authStore.user.phone_number || '';
    formData.gender = authStore.user.gender || '';
    formData.date_of_birth = authStore.user.date_of_birth || '';
    formData.bio = authStore.user.bio || '';
    formData.newsletter_subscribed = authStore.user.newsletter_subscribed || false;
  }
};

onMounted(() => {
  populateForm();
});

// Bersihkan memori URL pratinjau saat komponen dihancurkan
onUnmounted(() => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
  }
});

const triggerFileInput = () => {
  if (!isEditing.value) return;
  fileInput.value?.click();
};

const handleAvatarChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  // Validasi ukuran
  if (file.size > 2 * 1024 * 1024) {
    errorMessage.value = "Image size cannot exceed 2MB.";
    return;
  }

  // Buat pratinjau lokal
  selectedFile.value = file;
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value);
  previewUrl.value = URL.createObjectURL(file);
  errorMessage.value = '';
};

const cancelEdit = () => {
  populateForm();
  isEditing.value = false;
  errorMessage.value = '';
  
  // Reset state gambar
  selectedFile.value = null;
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
    previewUrl.value = null;
  }
};

const handleSave = async () => {
  isSaving.value = true;
  saveSuccess.value = false;
  errorMessage.value = '';

  try {
    // 1. Unggah gambar terlebih dahulu jika ada yang dipilih
    if (selectedFile.value) {
      await authStore.uploadAvatar(selectedFile.value);
    }

    // 2. Update data profil lainnya
    let formattedDate = formData.date_of_birth;
    
    // Jika formatnya YYYY-MM-DD (dari datepicker), tambahkan waktu dummy agar ISO valid
    // tanpa merubah tanggal aslinya.
    if (formattedDate && !formattedDate.includes('T')) {
       formattedDate = `${formattedDate}T00:00:00Z`;
    }

    await authStore.updateProfile({
      ...formData,
      date_of_birth: formattedDate || undefined
    });
    
    // Sinkronkan ulang form dengan data terbaru dari store setelah update sukses
    populateForm();
    
    saveSuccess.value = true;
    isEditing.value = false;
    
    // Reset state gambar setelah sukses
    selectedFile.value = null;
    if (previewUrl.value) {
      URL.revokeObjectURL(previewUrl.value);
      previewUrl.value = null;
    }

    setTimeout(() => { saveSuccess.value = false; }, 3000);
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || 'Failed to update profile. Please try again.';
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 md:p-8">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Profile Settings</h1>
        <p class="text-gray-500 mt-1">Manage your personal information and preferences.</p>
      </div>
      <BaseButton 
        v-if="!isEditing" 
        variant="outline" 
        class="md:w-auto px-6"
        @click="isEditing = true"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
        Edit Profile
      </BaseButton>
    </div>

    <!-- Alert Messages -->
    <BaseAlert 
      v-if="saveSuccess" 
      type="success" 
      message="Profile updated successfully." 
      class="mb-6"
    />

    <BaseAlert 
      v-if="errorMessage" 
      type="error" 
      :message="errorMessage" 
      class="mb-6"
    />

    <!-- Avatar Upload Section -->
    <div class="flex items-center gap-6 mb-10 pb-10 border-b border-gray-100">
      <div 
        class="relative group" 
        :class="isEditing ? 'cursor-pointer' : 'cursor-default'"
        @click="triggerFileInput"
      >
        <div class="w-24 h-24 rounded-full overflow-hidden border-2 border-gray-100 bg-gray-50 flex items-center justify-center">
          <!-- Gunakan previewUrl jika ada, jika tidak gunakan URL dari store -->
          <img 
            :src="previewUrl || authStore.user?.avatar_url || 'https://api.dicebear.com/7.x/initials/svg?seed=' + (authStore.user?.full_name || 'User')" 
            alt="User Avatar" 
            class="w-full h-full object-cover transition-opacity"
            :class="isEditing && 'group-hover:opacity-75'"
          />
        </div>
        <div 
          v-if="isEditing"
          class="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
        </div>
        <input 
          type="file" 
          ref="fileInput" 
          class="hidden" 
          accept="image/jpeg, image/png, image/webp"
          @change="handleAvatarChange"
        />
      </div>
      <div>
        <h3 class="font-semibold text-gray-900">Profile Picture</h3>
        <p class="text-sm text-gray-500 mt-1">PNG, JPG or WebP. Max 2MB.</p>
        <div v-if="isEditing" class="flex items-center gap-3 mt-2">
          <button 
            @click="triggerFileInput" 
            class="text-sm font-medium text-[#7b00ff] hover:text-[#6a00e0] bg-transparent border-none cursor-pointer p-0"
          >
            Change picture
          </button>
          <span v-if="selectedFile" class="text-xs text-green-600 font-medium flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 11"></polyline></svg>
            Ready to save
          </span>
        </div>
        <span v-else class="mt-2 inline-block text-xs bg-gray-100 text-gray-500 px-2 py-1 rounded-md">Click Edit to change</span>
      </div>
    </div>

    <!-- Edit Form -->
    <form @submit.prevent="handleSave" class="space-y-6">
      
      <!-- Grid 2 Columns -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <BaseInput 
          id="fullName" 
          label="Full Name" 
          type="text" 
          v-model="formData.full_name" 
          placeholder="John Doe" 
          :readonly="!isEditing"
          required 
        />
        <BaseInput 
          id="username" 
          label="Username" 
          type="text" 
          v-model="formData.username" 
          placeholder="johndoe" 
          :readonly="!isEditing"
        />
        <BaseInput 
          id="email" 
          label="Email Address" 
          type="email" 
          :model-value="authStore.user?.email" 
          disabled
          readonly
          class="opacity-60 cursor-not-allowed"
          title="Email cannot be changed directly"
        />
        <BaseInput 
          id="phone" 
          label="Phone Number" 
          type="tel" 
          v-model="formData.phone_number" 
          placeholder="+1234567890" 
          :readonly="!isEditing"
        />
        
        <BaseSelect
          id="gender"
          label="Gender"
          v-model="formData.gender"
          :options="genderOptions"
          placeholder="Select gender"
          :readonly="!isEditing"
        />

        <BaseDatePicker
          id="dob"
          label="Date of Birth"
          v-model="formData.date_of_birth"
          placeholder="Select your birthday"
          :readonly="!isEditing"
        />
      </div>

      <!-- Full Width Textarea -->
      <div class="flex flex-col gap-1.5">
        <label 
          for="bio" 
          class="text-sm font-semibold text-gray-700"
          :class="!isEditing && 'opacity-70'"
        >
          Bio
        </label>
        <textarea 
          id="bio" 
          v-model="formData.bio"
          rows="4"
          placeholder="Tell us a little about yourself..."
          :disabled="!isEditing"
          class="p-4 rounded-xl border bg-white text-gray-900 focus:outline-none focus:ring-2 focus:ring-[#7b00ff]/20 focus:border-[#7b00ff] transition-all hover:border-gray-400 w-full resize-none disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed"
          :class="errorMessage ? 'border-red-500' : 'border-gray-200'"
        ></textarea>
      </div>

      <!-- Toggles -->
      <div class="flex items-center gap-3 pt-2">
        <BaseToggle 
          id="newsletter"
          v-model="formData.newsletter_subscribed"
          label="Subscribe to newsletter and promos"
          :disabled="!isEditing"
        />
      </div>

      <!-- Form Actions -->
      <div v-if="isEditing" class="pt-6 border-t border-gray-100 flex flex-col md:flex-row justify-end gap-3">
        <BaseButton 
          type="button" 
          variant="outline"
          class="md:w-auto px-8"
          @click="cancelEdit"
          :disabled="isSaving"
        >
          Cancel
        </BaseButton>
        <BaseButton 
          type="submit" 
          :is-loading="isSaving" 
          loading-text="Saving..."
          class="md:w-auto px-8"
        >
          Save Changes
        </BaseButton>
      </div>

    </form>
  </div>
</template>
