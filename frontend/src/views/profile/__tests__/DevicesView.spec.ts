import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import DevicesView from '../DevicesView.vue';
import SessionCard from '../../../components/SessionCard.vue';
import { useAuthStore } from '../../../stores/auth';

describe('DevicesView.vue', () => {
  let wrapper: any;
  const mockSessions = [
    {
      id: '1',
      user_id: 'u1',
      ip_address: '127.0.0.1',
      user_agent: 'Mozilla/5.0 (Windows NT 10.0)',
      device_type: 'Desktop',
      location: 'Jakarta, ID',
      last_active_at: new Date().toISOString(),
      expires_at: new Date().toISOString(),
      is_current: true
    },
    {
      id: '2',
      user_id: 'u1',
      ip_address: '192.168.1.1',
      user_agent: 'Mozilla/5.0 (iPhone; CPU iPhone OS)',
      device_type: 'Mobile',
      location: 'Bandung, ID',
      last_active_at: new Date().toISOString(),
      expires_at: new Date().toISOString(),
      is_current: false
    }
  ];

  beforeEach(() => {
    wrapper = mount(DevicesView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: {
                activeSessions: mockSessions,
                isLoading: false
              }
            }
          })
        ]
      }
    });
  });

  it('renders the header and session list correctly', () => {
    expect(wrapper.find('h1').text()).toBe('Devices & Security');
    // Check if SessionCard components are rendered (should be 2 based on mock)
    const cards = wrapper.findAllComponents(SessionCard);
    expect(cards).toHaveLength(2);
  });

  it('shows the "Logout all other sessions" button when there are multiple sessions', () => {
    // When BaseButton is rendered, we look for the component or text
    const logoutOthersBtn = wrapper.findComponent({ name: 'BaseButton' }); 
    expect(logoutOthersBtn.exists()).toBe(true);
    expect(logoutOthersBtn.text()).toContain('Logout all other sessions');
  });

  it('calls fetchSessions on mount', () => {
    const authStore = useAuthStore();
    expect(authStore.fetchSessions).toHaveBeenCalledOnce();
  });

  it('shows empty state when no sessions exist', async () => {
    const authStore = useAuthStore();
    authStore.activeSessions = [];
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain('No active devices');
  });
});
