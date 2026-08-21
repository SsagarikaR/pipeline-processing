import { useCallback, useEffect, useState } from 'react';
import { Link, useParams, useNavigate } from 'react-router-dom';
import { jobService as api } from '../service/jobService';
import type { Result, JobError } from '../types/job';
import { RefreshCw } from 'lucide-react';
import StatusBadge from '../components/StatusBadge';
import usePolling from '../hooks/usePolling';
import AppButton from '../components/AppButton';
import ConfirmModal from '../components/ConfirmModal';

const TERMINAL_STATUSES = ['completed', 'failed', 'cancelled'] as const;

export default function JobDetail() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [results, setResults] = useState<Result[]>([]);
    const [errors, setErrors] = useState<JobError[]>([]);
    const [tab, setTab] = useState<'progress' | 'results' | 'errors'>('progress');

    const [previewData, setPreviewData] = useState<string | null>(null);
    const [previewLoading, setPreviewLoading] = useState(false);
    
    const [deleteModalOpen, setDeleteModalOpen] = useState(false);
    const [cancelModalOpen, setCancelModalOpen] = useState(false);

    // useParams can technically return undefined if the route param is missing —
    // guard here so every api.* call below can safely assume `id` is a string.
    if (!id) {
        return <div className="p-8 text-danger-600">Invalid job ID.</div>;
    }

    async function handlePreview(e: React.MouseEvent, url: string) {
        e.preventDefault();
        setPreviewLoading(true);
        setPreviewData(null);
        try {
            const res = await fetch(url);
            const text = await res.text();
            try {
                // Try to pretty-print JSON if possible
                const json = JSON.parse(text);
                setPreviewData(JSON.stringify(json, null, 2));
            } catch {
                setPreviewData(text);
            }
        } catch(err) {
            setPreviewData("Failed to load preview.");
        } finally {
            setPreviewLoading(false);
        }
    }

    const fetchProgress = useCallback(() => api.getProgress(id), [id]);
    const { data: progress, error: progressError, refresh: refreshProgress } = usePolling(fetchProgress, 15 * 60 * 1000, true);

    const isTerminal = progress ? (TERMINAL_STATUSES as readonly string[]).includes(progress.status) : false;

    const fetchResults = useCallback(() => {
        if (isTerminal) {
            api.getResults(id).then(res => setResults(res || [])).catch(() => { });
        }
    }, [isTerminal, id]);

    const fetchErrors = useCallback(() => {
        api.getErrors(id).then(errs => setErrors(errs || [])).catch(() => { });
    }, [id]);

    useEffect(() => {
        fetchResults();
    }, [fetchResults]);

    useEffect(() => {
        fetchErrors();
        const interval = setInterval(fetchErrors, 15 * 60 * 1000);
        return () => clearInterval(interval);
    }, [fetchErrors]);

    async function handleManualRefresh() {
        await refreshProgress();
        fetchResults();
        fetchErrors();
    }

    async function handleCancel() {
        if (id) {
            await api.cancelJob(id);
            setCancelModalOpen(false);
        }
    }

    async function handleDelete() {
        if (id) {
            await api.deleteJob(id);
            setDeleteModalOpen(false);
            navigate('/');
        }
    }

    if (progressError) {
        return <div className="p-8 text-danger-600">Failed to load job — check that it exists.</div>;
    }
    if (!progress) {
        return <div className="p-8 text-neutral-500">Loading…</div>;
    }

    return (
        <div className="max-w-3xl mx-auto p-6">
            <div className="mb-4">
                <Link to="/" className="text-sm text-brand-600 hover:underline flex items-center gap-1">
                    <span>←</span> Back to Jobs
                </Link>
            </div>
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-3">
                    <h1 className="text-2xl font-semibold text-neutral-900">Job #{progress.jobId}</h1>
                    <StatusBadge status={progress.status} />
                    <button 
                        onClick={handleManualRefresh} 
                        className="p-1.5 text-neutral-400 hover:text-neutral-900 hover:bg-neutral-100 rounded-lg transition-colors"
                        title="Refresh Job Data"
                    >
                        <RefreshCw size={16} />
                    </button>
                </div>
                <div className="flex gap-2">
                    {!isTerminal && (
                        <AppButton
                            variant="secondary"
                            onClick={() => setCancelModalOpen(true)}
                        >
                            Cancel
                        </AppButton>
                    )}
                    <AppButton
                        variant="danger"
                        onClick={() => setDeleteModalOpen(true)}
                    >
                        Delete
                    </AppButton>
                </div>
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
                <MetricCard label="Processed" value={progress.processed} />
                <MetricCard label="Errors" value={progress.errorCount} tone={progress.errorCount > 0 ? 'red' : 'default'} />
                <MetricCard label="Percent Complete" value={`${progress.percentComplete.toFixed(1)}%`} />
                <MetricCard label="Rate" value={progress.recordsPerSec ? `${progress.recordsPerSec.toFixed(1)}/s` : '—'} />
            </div>

            <div className="text-sm text-neutral-500 mb-6 flex gap-6">
                <span>Started: {progress.startedAt ? new Date(progress.startedAt).toLocaleString() : '—'}</span>
                <span>Completed: {progress.completedAt ? new Date(progress.completedAt).toLocaleString() : '—'}</span>
                {progress.exportUrl && (
                    <span>
                        Export: <a href={progress.exportUrl} target="_blank" rel="noreferrer" className="text-brand-600 hover:underline">{progress.exportUrl}</a>
                        <button onClick={(e) => handlePreview(e, progress.exportUrl!)} className="ml-3 text-xs bg-brand-50 text-brand-700 px-2 py-1 rounded hover:bg-brand-100">
                            Preview
                        </button>
                    </span>
                )}
            </div>

            {(previewData !== null || previewLoading) && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-lg shadow-xl w-full max-w-4xl flex flex-col max-h-[80vh]">
                        <div className="flex items-center justify-between p-4 border-b">
                            <h3 className="font-semibold">File Preview</h3>
                            <AppButton variant="ghost" size="sm" onClick={() => {setPreviewData(null); setPreviewLoading(false);}}>✕</AppButton>
                        </div>
                        <div className="p-4 overflow-auto bg-neutral-50 flex-1">
                            {previewLoading ? (
                                <div className="text-neutral-500 text-center py-8">Loading preview...</div>
                            ) : (
                                <pre className="text-xs text-neutral-800 font-mono whitespace-pre-wrap">{previewData}</pre>
                            )}
                        </div>
                    </div>
                </div>
            )}

            <div className="border-b border-neutral-200 mb-4">
                <nav className="flex gap-6">
                    {(['progress', 'results', 'errors'] as const).map((t) => (
                        <button
                            key={t}
                            onClick={() => setTab(t)}
                            className={`pb-3 text-sm font-medium border-b-2 transition-colors ${tab === t ? 'border-brand-600 text-brand-600' : 'border-transparent text-neutral-500 hover:text-neutral-700'
                                }`}
                        >
                            {t === 'results' ? `Results (${results?.length || 0})` : t === 'errors' ? `Errors (${errors?.length || 0})` : 'Live Progress'}
                        </button>
                    ))}
                </nav>
            </div>

            {tab === 'progress' && (
                <div className="text-sm text-neutral-600">
                    {isTerminal ? `Job finished with status "${progress.status}".` : 'Job is running — this page auto-refreshes every 2 seconds.'}
                </div>
            )}
            {tab === 'results' && <ResultsTable results={results} isTerminal={isTerminal} />}
            {tab === 'errors' && <ErrorsTable errors={errors} />}

            <ConfirmModal 
                isOpen={deleteModalOpen}
                title="Delete Job"
                message="Delete this job and its artifacts? This cannot be undone."
                onConfirm={handleDelete}
                onCancel={() => setDeleteModalOpen(false)}
            />
            <ConfirmModal 
                isOpen={cancelModalOpen}
                title="Cancel Job"
                message="Cancel this job?"
                onConfirm={handleCancel}
                onCancel={() => setCancelModalOpen(false)}
            />
        </div>
    );
}

interface MetricCardProps {
    label: string;
    value: string | number;
    tone?: 'default' | 'red';
}

function MetricCard({ label, value, tone = 'default' }: MetricCardProps) {
    const toneClass = tone === 'red' ? 'text-danger-600' : 'text-neutral-900';
    return (
        <div className="bg-white border border-neutral-200 rounded-lg p-3">
            <div className="text-xs text-neutral-500 mb-1">{label}</div>
            <div className={`text-lg font-semibold ${toneClass}`}>{value}</div>
        </div>
    );
}

interface ResultsTableProps {
    results: Result[];
    isTerminal: boolean;
}

function ResultsTable({ results, isTerminal }: ResultsTableProps) {
    if (!isTerminal) {
        return <div className="text-sm text-neutral-400 py-8 text-center">Results will appear once the job finishes.</div>;
    }
    if (!results || results.length === 0) {
        return <div className="text-sm text-neutral-400 py-8 text-center">No results produced.</div>;
    }
    return (
        <table className="w-full text-sm">
            <thead>
                <tr className="text-left text-neutral-500 border-b border-neutral-200">
                    <th className="py-2 font-medium">Group Key</th>
                    <th className="py-2 font-medium">Aggregated Value</th>
                </tr>
            </thead>
            <tbody>
                {results.map((r) => (
                    <tr key={r.id} className="border-b border-neutral-100">
                        <td className="py-2 font-mono text-xs text-neutral-700">{r.group_key}</td>
                        <td className="py-2 text-neutral-900">{r.aggregated_value.toLocaleString()}</td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
}

interface ErrorsTableProps {
    errors: JobError[];
}

function ErrorsTable({ errors }: ErrorsTableProps) {
    if (!errors || errors.length === 0) {
        return <div className="text-sm text-neutral-400 py-8 text-center">No errors reported.</div>;
    }
    return (
        <div className="space-y-2">
            {errors.map((e) => (
                <div key={e.id} className="bg-danger-50 border border-danger-200 rounded-lg p-3 text-sm">
                    <div className="flex items-center gap-2 mb-1">
                        <span className="px-2 py-0.5 bg-danger-100 text-danger-700 rounded text-xs font-medium">{e.stage}</span>
                        <span className="text-xs text-neutral-400">{new Date(e.created_at).toLocaleTimeString()}</span>
                    </div>
                    <div className="text-danger-800">{e.error_message}</div>
                    {e.record_data && <pre className="mt-1 text-xs text-neutral-500 overflow-x-auto">{e.record_data}</pre>}
                </div>
            ))}
        </div>
    );
}