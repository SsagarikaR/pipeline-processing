import { ExportConfig } from '../types/job';

import { ExportInputProps } from '../types/job';

export default function ExportInput({ exports, onChange }: ExportInputProps) {
  function update(i: number, fields: Partial<ExportConfig>) {
    const next = [...exports];
    next[i] = { ...next[i], ...fields };
    onChange(next);
  }

  function add() {
    onChange([...exports, { type: 's3', path: '' }]);
  }

  function remove(i: number) {
    onChange(exports.filter((_, idx) => idx !== i));
  }

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">Exports</label>
      {exports.map((exp, i) => (
        <div key={i} className="flex gap-2 items-center">
          <span className="text-sm text-neutral-500 bg-neutral-50 px-3 py-2 rounded-lg border border-neutral-200">S3</span>
          <input
            placeholder="e.g., results.json"
            value={exp.path}
            onChange={(e) => {
              update(i, { path: e.target.value, type: 's3' });
            }}
            className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
          />
          <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">✕</button>
        </div>
      ))}
      <button type="button" onClick={add} className="text-sm text-brand-600 hover:text-brand-700 font-medium">
        + Add export
      </button>
    </div>
  );
}