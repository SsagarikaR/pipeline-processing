// src/components/TransformInput.tsx
import type { TransformConfig, TransformName } from '../types/job';

interface TransformInputProps {
  transforms: TransformConfig[];
  onChange: (transforms: TransformConfig[]) => void;
}

export default function TransformInput({ transforms, onChange }: TransformInputProps) {
  function updateName(i: number, name: TransformName) {
    const next = [...transforms];
    next[i] = { ...next[i], name };
    onChange(next);
  }

  function updateParam(i: number, key: string, value: string) {
    const next = [...transforms];
    next[i] = { ...next[i], params: { ...next[i].params, [key]: value } };
    onChange(next);
  }

  function add() {
    onChange([...transforms, { name: 'uppercase', params: { field: '' } }]);
  }

  function remove(i: number) {
    onChange(transforms.filter((_, idx) => idx !== i));
  }

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">Transforms</label>
      {transforms.map((t, i) => (
        <div key={i} className="flex gap-2 items-center">
          <select
            value={t.name}
            onChange={(e) => updateName(i, e.target.value as TransformName)}
            className="border border-neutral-300 rounded-lg px-3 py-2 text-sm"
          >
            <option value="uppercase">uppercase</option>
            <option value="lowercase">lowercase</option>
          </select>

          <input
            placeholder="field name"
            value={t.params.field ?? ''}
            onChange={(e) => updateParam(i, 'field', e.target.value)}
            className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
          />
          <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">✕</button>
        </div>
      ))}
      <button type="button" onClick={add} className="text-sm text-brand-600 hover:text-brand-700 font-medium">
        + Add transform
      </button>
    </div>
  );
}