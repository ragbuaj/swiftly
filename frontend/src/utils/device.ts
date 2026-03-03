export const formatLastActive = (dateStr: string) => {
  const date = new Date(dateStr);
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);
  
  if (diffInSeconds < 60) return 'Just now';
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
  return date.toLocaleDateString(undefined, { day: 'numeric', month: 'short', year: 'numeric' });
};

export const getDeviceInfo = (ua: string) => {
  const os = ua.includes('Windows') ? 'Windows' : 
             ua.includes('Macintosh') ? 'macOS' : 
             ua.includes('Android') ? 'Android' : 
             ua.includes('iPhone') || ua.includes('iPad') ? 'iOS' : 'Linux';
             
  const browser = ua.includes('Chrome') ? 'Chrome' : 
                  ua.includes('Firefox') ? 'Firefox' : 
                  ua.includes('Safari') ? 'Safari' : 
                  ua.includes('Edg') ? 'Edge' : 'Browser';
                  
  return { os, browser };
};
