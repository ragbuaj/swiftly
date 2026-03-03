import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import ProfileSettings from '../../views/profile/ProfileSettings.vue';
import { useAuthStore } from '../../stores/auth';

describe('ProfileSettings.vue', () => {
  let wrapper: any;
  
  beforeEach(() => {
    // Mount component with Pinia testing plugin
    wrapper = mount(ProfileSettings, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: {
                user: {
                  id: '123',
                  email: 'test@example.com',
                  full_name: 'Test User',
                  username: 'testuser',
                  phone_number: '123456',
                  gender: 'male',
                  bio: 'A simple bio',
                  newsletter_subscribed: true,
                }
              }
            }
          })
        ]
      }
    });
  });

  it('renders correctly and populates form with user data', () => {
    const authStore = useAuthStore();
    
    // Check if store was populated
    expect(authStore.user?.full_name).toBe('Test User');
    
    // Check if form inputs are rendered and populated correctly
    const fullNameInput = wrapper.find('input#fullName');
    expect(fullNameInput.exists()).toBe(true);
    // Since we use BaseInput component, we might not get direct value access, 
    // but we can check if it rendered the correct attributes
    
    // Check if Bio is rendered
    const bioTextarea = wrapper.find('textarea#bio');
    expect(bioTextarea.exists()).toBe(true);
    expect(bioTextarea.element.value).toBe('A simple bio');
  });

  it('calls updateProfile when form is submitted', async () => {
    const authStore = useAuthStore();
    
    // Find form and trigger submit
    const form = wrapper.find('form');
    await form.trigger('submit.prevent');
    
    // Check if the store action was called
    expect(authStore.updateProfile).toHaveBeenCalledOnce();
  });
});
