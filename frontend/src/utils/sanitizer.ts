/**
 * Sanitizes input string by removing HTML tags and trimming whitespace.
 * @param {string} text 
 * @returns {string}
 */
export const sanitize = (text) => {
  if (typeof text !== 'string') return text;
  // Remove potential script tags and HTML
  return text.replace(/<[^>]*>?/gm, '').trim();
};

/**
 * Normalizes email by trimming and converting to lowercase.
 * @param {string} email 
 * @returns {string}
 */
export const normalizeEmail = (email) => {
  if (typeof email !== 'string') return '';
  return email.trim().toLowerCase();
};
