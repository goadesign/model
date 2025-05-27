/**
 * Robust platform detection utilities
 */

/**
 * Detects if the user is on a Mac platform using multiple detection methods
 * for maximum compatibility across browsers and future-proofing.
 */
export const isMac = (): boolean => {
  if (typeof navigator === 'undefined') return false;
  
  // Method 1: Check userAgentData (modern browsers, most reliable)
  if ('userAgentData' in navigator && (navigator as any).userAgentData) {
    const platform = (navigator as any).userAgentData.platform;
    if (platform && platform.toLowerCase().includes('mac')) {
      return true;
    }
  }
  
  // Method 2: Check userAgent string (widely supported)
  const userAgent = navigator.userAgent.toLowerCase();
  if (userAgent.includes('mac os') || userAgent.includes('macintosh')) {
    return true;
  }
  
  // Method 3: Check platform (fallback, deprecated but still widely supported)
  if (navigator.platform) {
    const platform = navigator.platform.toLowerCase();
    if (platform.includes('mac') || platform.includes('darwin')) {
      return true;
    }
  }
  
  // Method 4: Check for Mac-specific features as additional validation
  try {
    const testEvent = new KeyboardEvent('keydown', { metaKey: true });
    if (testEvent.metaKey !== undefined) {
      // Additional heuristic: Mac typically has different key layouts
      return /mac|darwin|os x/i.test(navigator.userAgent);
    }
  } catch (e) {
    // Ignore errors in older browsers
  }
  
  return false;
};

/**
 * Gets the appropriate modifier key name for the current platform
 */
export const getModifierKeyName = (): string => {
  return isMac() ? 'Cmd' : 'Ctrl';
};

/**
 * Gets the appropriate modifier key property for keyboard events
 */
export const getModifierKeyProperty = (event: KeyboardEvent): boolean => {
  return isMac() ? event.metaKey : event.ctrlKey;
};

/**
 * Gets the appropriate Alt key name for the current platform
 */
export const getAltKeyName = (): string => {
  return isMac() ? 'Option' : 'Alt';
}; 