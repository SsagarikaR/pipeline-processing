import type { SourceConfig, SourceType } from '../types/job';

interface SourceInputProps {
  sources: SourceConfig[];
  onChange: (sources: SourceConfig[]) => void;
}

export default function SourceInput({ sources, onChange }: SourceInputProps) {
  function update(i: number, field: keyof SourceConfig, value: string) {
    const next = [...sources];
    next[i] = { ...next[i], [field]: value };
    onChange(next);
  }

  function add() {
    onChange([...sources, { type: 'csv', path: '' }]);
  }

  function remove(i: number) {
    onChange(sources.filter((_, idx) => idx !== i));
  }

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">Sources</label>
      {sources.map((src, i) => (
        <div key={i} className="flex gap-2">
          <select
            value={src.type}
            onChange={(e) => update(i, 'type', e.target.value as SourceType)}
            className="border border-neutral-300 rounded-lg px-3 py-2 text-sm"
          >
            <option value="csv">CSV</option>
            <option value="json">JSON</option>
          </select>
          <div className="flex-1 flex gap-2">
            <input
              type="text"
              placeholder="/path/to/file or Data URI"
              value={src.path}
              onChange={(e) => update(i, 'path', e.target.value)}
              className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            />
            <label className="cursor-pointer border border-neutral-300 rounded-lg px-3 py-2 text-sm bg-neutral-50 hover:bg-neutral-100 flex items-center shrink-0">
              <span>Upload</span>
              <input
                type="file"
                className="hidden"
                accept={src.type === 'json' ? '.json' : '.csv'}
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    const reader = new FileReader();
                    reader.onload = (event) => {
                      if (typeof event.target?.result === 'string') {
                        update(i, 'path', event.target.result);
                      }
                    };
                    reader.readAsDataURL(file);
                  }
                  e.target.value = ''; // reset
                }}
              />
            </label>
          </div>
          <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">
            ✕
          </button>
        </div>
      ))}
      <button type="button" onClick={add} className="text-sm text-brand-600 hover:text-brand-700 font-medium">
        + Add source
      </button>
    </div>
  );
}