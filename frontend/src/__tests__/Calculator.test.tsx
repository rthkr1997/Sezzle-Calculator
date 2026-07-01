import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Calculator } from '../components/Calculator';

const fetchMock = vi.fn();
vi.stubGlobal('fetch', fetchMock);

describe('Calculator', () => {
  beforeEach(() => fetchMock.mockReset());

  it('submits a chained expression and renders the result', async () => {
    fetchMock.mockResolvedValue({ ok: true, json: async () => ({ result: 22 }) });
    render(<Calculator />);
    await userEvent.type(screen.getByLabelText(/expression/i), '12 + 4 * (8 - 3) / 2');
    await userEvent.click(screen.getByRole('button', { name: /calculate/i }));
    expect(await screen.findByText(/Answer:\s*22/)).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:8080/api/calculate', expect.objectContaining({ method: 'POST' }));
  });

  it('shows backend validation errors', async () => {
    fetchMock.mockResolvedValue({ ok: false, json: async () => ({ error: 'division by zero' }) });
    render(<Calculator />);
    await userEvent.type(screen.getByLabelText(/expression/i), '10/0');
    await userEvent.click(screen.getByRole('button', { name: /calculate/i }));
    expect(await screen.findByRole('alert')).toHaveTextContent('division by zero');
  });
});
