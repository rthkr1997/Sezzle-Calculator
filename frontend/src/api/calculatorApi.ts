export type CalculateResponse = { result: number; expression?: string; operation?: string };
export type ApiError = { error: string };

const baseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

export async function calculateExpression(expression: string): Promise<CalculateResponse> {
  const response = await fetch(`${baseUrl}/api/calculate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ expression })
  });
  const payload = await response.json();
  if (!response.ok) {
    throw new Error((payload as ApiError).error ?? 'Calculation failed');
  }
  return payload as CalculateResponse;
}
