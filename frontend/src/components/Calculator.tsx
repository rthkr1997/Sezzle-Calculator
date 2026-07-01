import { FormEvent, useMemo, useState } from 'react';
import { calculateExpression } from '../api/calculatorApi';

const buttons = ['7','8','9','/','sqrt(','4','5','6','*','^','1','2','3','-','%','0','.','(',')','+'];

function isPotentiallyValidExpression(value: string): boolean {
  return /^[0-9+\-*/^%().\sA-Za-z]*$/.test(value);
}

export function Calculator() {
  const [expression, setExpression] = useState('');
  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState('');
  const [history, setHistory] = useState<string[]>([]);
  const canSubmit = useMemo(() => expression.trim().length > 0 && !error, [expression, error]);

  const updateExpression = (next: string) => {
    setResult(null);
    if (!isPotentiallyValidExpression(next)) {
      setError('Only numbers, arithmetic operators, parentheses, %, and sqrt are supported.');
      return;
    }
    setError('');
    setExpression(next);
  };

  const append = (value: string) => updateExpression(expression + value);

  const submit = async (event?: FormEvent) => {
    event?.preventDefault();
    if (!expression.trim()) {
      setError('Enter an expression to calculate.');
      return;
    }
    try {
      setError('');
      const response = await calculateExpression(expression);
      setResult(response.result);
      setHistory((items) => [`${expression} = ${response.result}`, ...items].slice(0, 5));
    } catch (caught) {
      setResult(null);
      setError(caught instanceof Error ? caught.message : 'Something went wrong.');
    }
  };

  return (
    <main className="calculator-shell">
      <section className="calculator-card" aria-label="Calculator">
        <div className="hero">
          <p className="eyebrow">Full-stack calculator</p>
          <h1>Expression Calculator</h1>
          <p>Use chained operations, parentheses, exponentiation, square roots, and percentages.</p>
        </div>
        <form onSubmit={submit} className="display-panel">
          <label htmlFor="expression">Expression</label>
          <input id="expression" value={expression} onChange={(e) => updateExpression(e.target.value)} placeholder="12 + 4 * (8 - 3) / 2" autoComplete="off" />
          {result !== null && <output className="result" aria-live="polite">{result}</output>}
          {error && <p className="error" role="alert">{error}</p>}
          <button className="primary" type="submit" disabled={!canSubmit}>Calculate</button>
        </form>
        <div className="keypad" aria-label="Calculator keypad">
          {buttons.map((button) => <button key={button} type="button" onClick={() => append(button)}>{button}</button>)}
          <button type="button" onClick={() => updateExpression(expression.slice(0, -1))}>⌫</button>
          <button type="button" onClick={() => { setExpression(''); setResult(null); setError(''); }}>C</button>
        </div>
        {history.length > 0 && <aside className="history"><h2>Recent calculations</h2>{history.map((item) => <p key={item}>{item}</p>)}</aside>}
      </section>
    </main>
  );
}
