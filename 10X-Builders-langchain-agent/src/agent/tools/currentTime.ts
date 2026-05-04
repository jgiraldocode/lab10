import { tool } from '@langchain/core/tools';
import { z } from 'zod';
import { BOGOTA_TZ, nowBogotaTimeHHMMSS, todayBogotaYYYYMMDD } from '../../lib/bogotaClock.js';

export const currentTimeTool = tool(
  async () => {
    const time = nowBogotaTimeHHMMSS();
    const date = todayBogotaYYYYMMDD();
    return (
      `Hora en Colombia (${BOGOTA_TZ}): ${time}. ` +
      `Fecha de hoy (para vuelos y calendario, YYYY-MM-DD): ${date}.`
    );
  },
  {
    name: 'current_time',
    description:
      'Devuelve la hora actual y la fecha de hoy en Colombia (zona ' +
      BOGOTA_TZ +
      '), incluida la fecha en formato YYYY-MM-DD. Úsala antes de flight_search cuando necesites saber qué día es hoy o construir fechas futuras.',
    schema: z.object({}),
  },
);
