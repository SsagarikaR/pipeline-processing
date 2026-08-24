import React from 'react';
import type { JobStatus, Result, JobError } from './job';

export interface AppButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
}

export type AppInputProps = React.InputHTMLAttributes<HTMLInputElement> & {
  /**
   * Visible, properly-associated label rendered above the input. Omit
   * it (and pass `aria-label` instead) for compact repeated-row fields
   * where a visible label per row would be redundant clutter.
   */
  label?: React.ReactNode;
};

export interface ConfirmModalProps {
  isOpen: boolean;
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export interface CreateJobModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface StatusBadgeProps {
  status: JobStatus;
}

export interface ErrorBoundaryProps {
  children: React.ReactNode;
}

export interface ErrorBoundaryState {
  hasError: boolean;
}

export interface MetricCardProps {
  label: string;
  value: string | number;
  tone?: 'default' | 'red';
}

export interface ResultsTableProps {
  results: Result[];
  isTerminal: boolean;
}

export interface ErrorsTableProps {
  errors: JobError[];
}
