export const stripHtml = (html: string): string =>
  typeof document !== 'undefined'
    ? (() => {
        const tmp = document.createElement('DIV');
        tmp.innerHTML = html;
        return tmp.textContent || tmp.innerText || '';
      })()
    : html;
