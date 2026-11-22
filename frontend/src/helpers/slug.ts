const TEXT_SLUG_SEPARATOR = '-';
const ITEM_SLUG_SEPARATOR = '-';

const charMap: Record<string, string> = {
  ł: 'l',
  Ł: 'l',
  ą: 'a',
  Ą: 'a',
  ć: 'c',
  Ć: 'c',
  ę: 'e',
  Ę: 'e',
  ń: 'n',
  Ń: 'n',
  ó: 'o',
  Ó: 'o',
  ś: 's',
  Ś: 's',
  ż: 'z',
  Ż: 'z',
  ź: 'z',
  Ź: 'z',
  ß: 'ss',
  ø: 'o',
  Ø: 'o',
  å: 'a',
  Å: 'a',
  æ: 'ae',
  Æ: 'ae',
  œ: 'oe',
  Œ: 'oe',
};

const normalizeText = (text: string): string => {
  return text
    .split('')
    .map((char) => charMap[char] ?? char)
    .join('')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
};

export const slugify = (text: string): string => {
  return normalizeText(text)
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, TEXT_SLUG_SEPARATOR)
    .replace(new RegExp(`${TEXT_SLUG_SEPARATOR}+`, 'g'), TEXT_SLUG_SEPARATOR)
    .replace(new RegExp(`^${TEXT_SLUG_SEPARATOR}|${TEXT_SLUG_SEPARATOR}$`, 'g'), '');
};

export const slugifyId = (id: number, text: string): string => {
  return `${id}${ITEM_SLUG_SEPARATOR}${slugify(text)}`;
};

export const unslugifyId = (slug: string): { id: number; name: string } => {
  const parts = slug.split(ITEM_SLUG_SEPARATOR);
  const id = parseInt(parts[0]);
  const name = parts.slice(1).join(TEXT_SLUG_SEPARATOR);
  return { id, name };
};
