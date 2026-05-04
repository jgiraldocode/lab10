import { tool } from '@langchain/core/tools';
import { z } from 'zod';
import { getEnv } from '../../config/env.js';
import { BOGOTA_TZ, todayBogotaYYYYMMDD } from '../../lib/bogotaClock.js';

const SERPAPI_SEARCH = 'https://serpapi.com/search';
const MAX_OPTIONS = 5;
const DATE_RE = /^\d{4}-\d{2}-\d{2}$/;

function normalizeLocationId(raw: string): string {
  const s = raw.trim();
  if (/^[a-zA-Z]{3}$/.test(s)) {
    return s.toUpperCase();
  }
  return s;
}

function normalizeOptionalCountryLang(raw: string | undefined): string | undefined {
  if (raw === undefined || raw === null) {
    return undefined;
  }
  const t = String(raw).trim().toLowerCase();
  if (t === '' || t === 'undefined') {
    return undefined;
  }
  return t.length === 2 ? t : undefined;
}

function normalizeCurrency(raw: string | undefined): string {
  if (raw === undefined || raw === null) {
    return 'USD';
  }
  const t = String(raw).trim().toUpperCase();
  if (t === '' || t === 'UNDEFINED') {
    return 'USD';
  }
  return t;
}

async function readSerpApiHttpError(res: Response): Promise<string> {
  const raw = await res.text();
  try {
    const j = JSON.parse(raw) as { error?: string };
    if (typeof j.error === 'string' && j.error.trim()) {
      return j.error.trim();
    }
  } catch {
    /* ignore */
  }
  const snippet = raw.replace(/\s+/g, ' ').trim().slice(0, 280);
  return snippet || res.statusText;
}

interface FlightLeg {
  airline?: string;
  flight_number?: string;
  departure_airport?: { id?: string; name?: string; time?: string };
  arrival_airport?: { id?: string; name?: string; time?: string };
}

interface BestFlightEntry {
  price?: number;
  type?: string;
  total_duration?: number;
  flights?: FlightLeg[];
  layovers?: { name?: string; duration?: number }[];
}

interface SerpFlightsPayload {
  search_metadata?: { status?: string };
  error?: string;
  best_flights?: BestFlightEntry[];
  other_flights?: BestFlightEntry[];
  search_parameters?: Record<string, unknown>;
}

function formatDurationMinutes(totalMinutes: number | undefined): string {
  if (totalMinutes === undefined || Number.isNaN(totalMinutes)) {
    return '?';
  }
  const h = Math.floor(totalMinutes / 60);
  const m = totalMinutes % 60;
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}

function summarizeFlightOption(entry: BestFlightEntry, index: number, currency: string): string {
  const price = entry.price ?? '?';
  const tripType = entry.type ?? 'vuelo';
  const duration = formatDurationMinutes(entry.total_duration);
  const segments = entry.flights ?? [];
  const airlines = [...new Set(segments.map((s) => s.airline).filter(Boolean))];
  const airlineStr = airlines.length ? airlines.join(' → ') : 'varias aerolíneas';
  const stops = entry.layovers?.length ?? Math.max(0, segments.length - 1);
  const routePreview = segments
    .map((s) => {
      const from = s.departure_airport?.id ?? '';
      const to = s.arrival_airport?.id ?? '';
      return from && to ? `${from}→${to}` : '';
    })
    .filter(Boolean)
    .join(', ');

  const lines = [
    `${index + 1}. Precio: ${price} ${currency} (${tripType})`,
    `   Duración total: ${duration}. Aerolíneas: ${airlineStr}. Escalas: ${stops}.`,
  ];
  if (routePreview) {
    lines.push(`   Tramos: ${routePreview}`);
  }
  return lines.join('\n');
}

function formatFlightsResponse(data: SerpFlightsPayload, currency: string): string {
  const status = data.search_metadata?.status;
  if (status && status !== 'Success') {
    return `La búsqueda en SerpApi no terminó correctamente (estado: ${status}).`;
  }
  if (data.error) {
    return `Error de SerpApi: ${data.error}`;
  }

  const primary = data.best_flights?.length ? data.best_flights : data.other_flights;
  if (!primary?.length) {
    return 'No se encontraron vuelos para estos parámetros. Prueba otras fechas o aeropuertos.';
  }

  const header =
    `Opciones encontradas (mostrando hasta ${MAX_OPTIONS}):\n` +
    `(Moneda: ${currency})\n\n`;
  const body = primary
    .slice(0, MAX_OPTIONS)
    .map((entry, i) => summarizeFlightOption(entry, i, currency))
    .join('\n\n');

  return header + body;
}

export const flightSearchTool = tool(
  async ({
    departure_id,
    arrival_id,
    outbound_date,
    trip_type,
    return_date,
    currency,
    gl,
    hl,
    adults,
  }) => {
    const env = getEnv();
    const apiKey = env.SERPAPI_API_KEY?.trim();
    if (!apiKey) {
      return (
        'SerpApi no está configurado: añade SERPAPI_API_KEY en env.local ' +
        '(https://serpapi.com).'
      );
    }

    if (trip_type === 'round_trip' && !return_date?.trim()) {
      return 'Para ida y vuelta debes indicar return_date (YYYY-MM-DD).';
    }

    const dep = normalizeLocationId(departure_id);
    const arr = normalizeLocationId(arrival_id);
    const outDate = outbound_date.trim();
    const retDate = return_date?.trim() ?? '';
    const cur = normalizeCurrency(currency);
    const glNorm = normalizeOptionalCountryLang(gl);
    const hlNorm = normalizeOptionalCountryLang(hl);

    if (!dep || !arr) {
      return 'departure_id y arrival_id no pueden estar vacíos (usa códigos IATA de 3 letras, ej. BOG, LHR).';
    }
    if (!DATE_RE.test(outDate)) {
      return `outbound_date debe ser YYYY-MM-DD; recibí: "${outbound_date}".`;
    }
    if (trip_type === 'round_trip') {
      if (!DATE_RE.test(retDate)) {
        return `return_date debe ser YYYY-MM-DD para ida y vuelta; recibí: "${return_date}".`;
      }
    }

    const today = todayBogotaYYYYMMDD();
    if (outDate < today) {
      return (
        `La fecha de salida (${outDate}) es anterior a hoy en Colombia (${today}, ${BOGOTA_TZ}). ` +
        'SerpApi no acepta fechas pasadas: llama primero a current_time para ver la fecha de hoy, ' +
        'luego usa esa fecha o una futura (YYYY-MM-DD).'
      );
    }
    if (trip_type === 'round_trip') {
      if (retDate < today) {
        return (
          `La fecha de regreso (${retDate}) es anterior a hoy en Colombia (${today}, ${BOGOTA_TZ}). ` +
          'Usa current_time para la fecha de hoy o una fecha futura.'
        );
      }
      if (retDate < outDate) {
        return (
          `return_date (${retDate}) no puede ser anterior a outbound_date (${outDate}).`
        );
      }
    }

    const params = new URLSearchParams({
      engine: 'google_flights',
      api_key: apiKey,
      departure_id: dep,
      arrival_id: arr,
      outbound_date: outDate,
      type: trip_type === 'one_way' ? '2' : '1',
      currency: cur,
      adults: String(
        typeof adults === 'number' && Number.isFinite(adults) && adults >= 1
          ? Math.floor(adults)
          : 1,
      ),
    });

    if (glNorm) {
      params.set('gl', glNorm);
    }
    if (hlNorm) {
      params.set('hl', hlNorm);
    }
    if (trip_type === 'round_trip' && retDate) {
      params.set('return_date', retDate);
    }

    const url = `${SERPAPI_SEARCH}?${params.toString()}`;

    let res: Response;
    try {
      res = await fetch(url);
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      return `No se pudo contactar a SerpApi: ${msg}`;
    }

    if (!res.ok) {
      const detail = await readSerpApiHttpError(res);
      return `SerpApi respondió HTTP ${res.status}: ${detail}`;
    }

    let data: SerpFlightsPayload;
    try {
      data = (await res.json()) as SerpFlightsPayload;
    } catch {
      return 'La respuesta de SerpApi no es JSON válido.';
    }

    return formatFlightsResponse(data, cur);
  },
  {
    name: 'flight_search',
    description:
      'Busca vuelos y precios con Google Flights vía SerpApi. ' +
      'Usa códigos IATA de aeropuerto (ej. BOG, LHR, MAD) o ids de ubicación soportados por SerpApi. ' +
      'outbound_date y return_date deben ser hoy o futuras (YYYY-MM-DD); nunca uses fechas pasadas. ' +
      'Si el usuario solo nombra ciudad/país, infiere el aeropuerto principal o pide aclaración.',
    schema: z.object({
      departure_id: z
        .string()
        .describe('Origen: código IATA o id de ubicación (ej. BOG, CDG).'),
      arrival_id: z
        .string()
        .describe('Destino: código IATA o id de ubicación (ej. LHR, MAD).'),
      outbound_date: z
        .string()
        .describe(
          'Fecha de salida YYYY-MM-DD: debe ser hoy o una fecha futura (no pasada).',
        ),
      trip_type: z
        .enum(['one_way', 'round_trip'])
        .describe('Ida simple o ida y vuelta.'),
      return_date: z
        .string()
        .optional()
        .describe(
          'Fecha de regreso YYYY-MM-DD (obligatoria en ida y vuelta): hoy o futura, posterior o igual a outbound_date.',
        ),
      currency: z
        .string()
        .default('USD')
        .describe('Moneda ISO para los precios (ej. USD, EUR, COP).'),
      gl: z
        .string()
        .optional()
        .describe(
          'País de Google: exactamente dos letras ISO (minúsculas), ej. us, co, uk. No uses el nombre del país.',
        ),
      hl: z
        .string()
        .optional()
        .describe(
          'Idioma de la búsqueda: exactamente dos letras, ej. es, en. No uses el nombre del idioma.',
        ),
      adults: z
        .number()
        .int()
        .min(1)
        .optional()
        .describe('Número de adultos (por defecto 1 en SerpApi si no se envía).'),
    }),
  },
);
