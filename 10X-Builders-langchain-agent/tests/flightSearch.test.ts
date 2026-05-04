import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flightSearchTool } from '../src/agent/tools/flightSearch.js';

describe('flightSearchTool', () => {
  const originalFetch = globalThis.fetch;
  const originalSerpKey = process.env.SERPAPI_API_KEY;

  beforeEach(() => {
    process.env.OPENROUTER_API_KEY = 'test-openrouter';
    process.env.SERPAPI_API_KEY = 'test-serp';
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
    if (originalSerpKey === undefined) {
      delete process.env.SERPAPI_API_KEY;
    } else {
      process.env.SERPAPI_API_KEY = originalSerpKey;
    }
    vi.restoreAllMocks();
  });

  it('formatea vuelos cuando SerpApi responde con best_flights', async () => {
    const payload = {
      search_metadata: { status: 'Success' },
      best_flights: [
        {
          price: 500,
          type: 'Round trip',
          total_duration: 480,
          flights: [
            {
              airline: 'IB',
              departure_airport: { id: 'MAD', time: '2026-05-03 10:00' },
              arrival_airport: { id: 'LHR', time: '2026-05-03 12:00' },
            },
          ],
          layovers: [],
        },
      ],
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      statusText: 'OK',
      json: async () => payload,
    } as Response);

    const result = await flightSearchTool.invoke({
      departure_id: 'MAD',
      arrival_id: 'LHR',
      outbound_date: '2026-05-10',
      trip_type: 'one_way',
      currency: 'EUR',
    });

    expect(result).toContain('500');
    expect(result).toContain('EUR');
    expect(result).toContain('IB');
    expect(globalThis.fetch).toHaveBeenCalled();
    const url = String(vi.mocked(globalThis.fetch).mock.calls[0]?.[0]);
    expect(url).toContain('engine=google_flights');
    expect(url).toContain('type=2');
    expect(url).toContain('adults=1');
  });

  it('incluye el mensaje JSON de SerpApi cuando HTTP no es OK', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      statusText: 'Bad Request',
      text: async () =>
        JSON.stringify({ error: 'Missing departure_id parameter.' }),
    } as Response);

    const result = await flightSearchTool.invoke({
      departure_id: 'BOG',
      arrival_id: 'LHR',
      outbound_date: '2026-06-01',
      trip_type: 'one_way',
      currency: 'USD',
    });

    expect(result).toContain('HTTP 400');
    expect(result).toContain('Missing departure_id');
  });

  it('normaliza códigos IATA a mayúsculas en la URL', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      statusText: 'OK',
      json: async () => ({
        search_metadata: { status: 'Success' },
        best_flights: [],
      }),
    } as Response);

    await flightSearchTool.invoke({
      departure_id: 'bog',
      arrival_id: 'lhr',
      outbound_date: '2026-06-01',
      trip_type: 'one_way',
      currency: 'USD',
    });

    const url = String(vi.mocked(globalThis.fetch).mock.calls[0]?.[0]);
    expect(url).toContain('departure_id=BOG');
    expect(url).toContain('arrival_id=LHR');
  });

  it('avisa si falta SERPAPI_API_KEY', async () => {
    delete process.env.SERPAPI_API_KEY;

    const result = await flightSearchTool.invoke({
      departure_id: 'BOG',
      arrival_id: 'LHR',
      outbound_date: '2026-06-01',
      trip_type: 'one_way',
      currency: 'USD',
    });

    expect(result).toContain('SERPAPI_API_KEY');
    expect(result).toContain('env.local');
  });

  it('exige return_date en ida y vuelta', async () => {
    const result = await flightSearchTool.invoke({
      departure_id: 'BOG',
      arrival_id: 'MAD',
      outbound_date: '2026-06-01',
      trip_type: 'round_trip',
      currency: 'USD',
    });

    expect(result).toContain('return_date');
  });

  it('rechaza outbound_date en el pasado sin llamar a SerpApi', async () => {
    const fetchSpy = vi.fn();
    globalThis.fetch = fetchSpy;

    const result = await flightSearchTool.invoke({
      departure_id: 'BOG',
      arrival_id: 'LHR',
      outbound_date: '2000-01-01',
      trip_type: 'one_way',
      currency: 'USD',
    });

    expect(result).toContain('anterior a hoy');
    expect(result).toContain('2000-01-01');
    expect(fetchSpy).not.toHaveBeenCalled();
  });

  it('rechaza return_date anterior a outbound en ida y vuelta', async () => {
    const fetchSpy = vi.fn();
    globalThis.fetch = fetchSpy;

    const result = await flightSearchTool.invoke({
      departure_id: 'BOG',
      arrival_id: 'MAD',
      outbound_date: '2099-12-10',
      trip_type: 'round_trip',
      return_date: '2099-12-01',
      currency: 'USD',
    });

    expect(result).toContain('no puede ser anterior');
    expect(fetchSpy).not.toHaveBeenCalled();
  });
});
