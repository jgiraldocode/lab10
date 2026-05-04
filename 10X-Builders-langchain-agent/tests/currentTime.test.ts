import { describe, expect, it } from "vitest";
import { currentTimeTool } from "../src/agent/tools/currentTime.js";

describe('currentTimeTool', () => {
  it('devuelve hora y fecha YYYY-MM-DD (Bogotá)', async () => {
    const result = await currentTimeTool.invoke({});
    expect(result).toMatch(/\d{1,2}:\d{2}:\d{2}/);
    expect(result).toMatch(/\d{4}-\d{2}-\d{2}/);
    expect(result).toContain('America/Bogota');
  });
});
