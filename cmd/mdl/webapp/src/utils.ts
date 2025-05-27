// Helper functions for the application

export function removeEmptyProps(obj: any) {
  return JSON.parse(JSON.stringify(obj));
}

export function camelToWords(camel: string) {
  const split = camel.replace(/([A-Z])/g, " $1");
  return split.charAt(0).toUpperCase() + split.slice(1);
}

export function getCurrentViewID() {
  const params = new URLSearchParams(document.location.search);
  return params.get('id') || '';
} 