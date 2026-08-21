import { AggregationConfig } from '../types/job';

interface Props {
  aggregations: AggregationConfig[];
  onChange: (aggs: AggregationConfig[]) => void;
}

const FIELD_TYPES = ['Number', 'String', 'Boolean', 'Date'];

const AGG_OPS: Record<string, string[]> = {
  Number: ['sum', 'avg', 'min', 'max', 'count'],
  String: ['count'],
  Boolean: ['count'],
  Date: ['min', 'max', 'count'],
};

export default function AggregationInput({ aggregations, onChange }: Props) {
  function update(i: number, field: keyof AggregationConfig | 'type', value: string) {
    const next: any[] = [...aggregations];
    
     next[i] = { ...next[i], [field]: value };
    
      if (field === 'type') {
      const availableOps = AGG_OPS[value] || [];
      if (!availableOps.includes(next[i].op)) {
        next[i].op = availableOps[0] || 'count';
      }
    }
    
    onChange(next);
  }

  function add() {
   onChange([...aggregations, { field: '', op: 'sum', groupBy: '', type: 'Number' } as any]);
  }

  function remove(i: number) {
    onChange(aggregations.filter((_, idx) => idx !== i));
  }

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">Aggregations</label>
      {aggregations.map((a, i) => {
        const currentType = (a as any).type || 'Number';
        const availableOps = AGG_OPS[currentType] || AGG_OPS.Number;
        
        return (
          <div key={i} className="flex gap-2">
            <select
              value={currentType}
              onChange={(e) => update(i, 'type', e.target.value)}
              className="border border-neutral-300 rounded-lg px-3 py-2 text-sm w-32"
            >
              {FIELD_TYPES.map(t => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
            <select
              value={a.op}
              onChange={(e) => update(i, 'op', e.target.value)}
              className="border border-neutral-300 rounded-lg px-3 py-2 text-sm w-28"
            >
              {availableOps.map(op => (
                <option key={op} value={op}>{op}</option>
              ))}
            </select>
            <input
              placeholder="field"
              value={a.field}
              onChange={(e) => update(i, 'field', e.target.value)}
              className="border border-neutral-300 rounded-lg px-3 py-2 text-sm w-32"
            />
            <input
              placeholder="group by (optional)"
              value={a.groupBy || ''}
              onChange={(e) => update(i, 'groupBy', e.target.value)}
              className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            />
            <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">✕</button>
          </div>
        );
      })}
      <button type="button" onClick={add} className="text-sm text-brand-600 hover:text-brand-700 font-medium">
        + Add aggregation
      </button>
    </div>
  );
}