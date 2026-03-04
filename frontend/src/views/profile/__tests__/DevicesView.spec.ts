import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import DevicesView from '../DevicesView.vue';
import SessionCard from '../../../components/SessionCard.vue';
import BaseModal from '../../../components/BaseModal.vue';
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
        ],
        stubs: {
          BaseModal: true,
          Teleport: true
        }
      }
    });
  });

  it('renders the header and session list correctly', () => {
    expect(wrapper.find('h1').text()).toBe('Devices & Security');
    const cards = wrapper.findAllComponents(SessionCard);
    expect(cards).toHaveLength(2);
  });

  it('shows the "Logout all other sessions" button when there are multiple sessions', () => {
    const buttons = wrapper.findAllComponents({ name: 'BaseButton' });
    const logoutOthersBtn = buttons.find((btn: any) => btn.text().includes('Logout all other sessions'));
    expect(logoutOthersBtn?.exists()).toBe(true);
  });

  it('opens logout modal when a session revoke is triggered', async () => {
    const card = wrapper.findComponent(SessionCard);
    await card.vm.$emit('revoke', '2');

    const modal = wrapper.findComponent(BaseModal);
    expect(modal.props('open')).toBe(true);
    expect(modal.props('title')).toBe('Logout device?');
  });

  it('opens logout all others modal when header button is clicked', async () => {
    const buttons = wrapper.findAllComponents({ name: 'BaseButton' });
    const logoutOthersBtn = buttons.find((btn: any) => btn.text().includes('Logout all other sessions'));
    
    await logoutOthersBtn.trigger('click');

    const modal = wrapper.findComponent(BaseModal);
    expect(modal.props('open')).toBe(true);
    expect(modal.props('title')).toBe('Logout all other devices?');
  });

  it('calls authStore.revokeSession when modal confirm is emitted', async () => {
    const authStore = useAuthStore();
    const card = wrapper.findComponent(SessionCard);
    await card.vm.$emit('revoke', '2');

    const modal = wrapper.findComponent(BaseModal);
    await modal.vm.$emit('confirm');

    expect(authStore.revokeSession).toHaveBeenCalledWith('2');
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
