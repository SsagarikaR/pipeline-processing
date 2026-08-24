import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import ConfirmModal from '../ConfirmModal';

describe('ConfirmModal', () => {
  it('does not render when isOpen is false', () => {
    render(
      <ConfirmModal 
        isOpen={false} 
        title="Test" 
        message="Msg" 
        onConfirm={vi.fn()} 
        onCancel={vi.fn()} 
      />
    );
    expect(screen.queryByText('Test')).not.toBeInTheDocument();
  });

  it('renders correctly when isOpen is true', () => {
    render(
      <ConfirmModal 
        isOpen={true} 
        title="Delete Job" 
        message="Are you sure?" 
        onConfirm={vi.fn()} 
        onCancel={vi.fn()} 
      />
    );
    expect(screen.getByText('Delete Job')).toBeInTheDocument();
    expect(screen.getByText('Are you sure?')).toBeInTheDocument();
  });

  it('calls onConfirm when confirm button is clicked', () => {
    const onConfirm = vi.fn();
    render(
      <ConfirmModal 
        isOpen={true} 
        title="Test" 
        message="Msg" 
        onConfirm={onConfirm} 
        onCancel={vi.fn()} 
      />
    );
    fireEvent.click(screen.getByText('Confirm'));
    expect(onConfirm).toHaveBeenCalledTimes(1);
  });
});
