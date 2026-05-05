export const VALID_CURRENCIES = ['USD', 'BAM', 'EURO'] as const;
export type Currency = (typeof VALID_CURRENCIES)[number];

export const CURRENCY_SYMBOLS: Record<Currency, string> = {
  USD: '$',
  BAM: 'KM',
  EURO: '€',
};

export const isCurrencyValid = (currency: string): currency is Currency => {
  return VALID_CURRENCIES.includes(currency as Currency);
};
