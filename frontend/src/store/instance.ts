import { create } from 'zustand';
import type { ClawInstance } from '@/types';

interface InstanceStore {
  instances: ClawInstance[];
  selectedInstance: ClawInstance | null;
  loading: boolean;
  error: string | null;
  setInstances: (instances: ClawInstance[]) => void;
  setSelectedInstance: (instance: ClawInstance | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  addInstance: (instance: ClawInstance) => void;
  updateInstance: (id: string, instance: Partial<ClawInstance>) => void;
  removeInstance: (id: string) => void;
}

export const useInstanceStore = create<InstanceStore>((set) => ({
  instances: [],
  selectedInstance: null,
  loading: false,
  error: null,

  setInstances: (instances) => set({ instances }),

  setSelectedInstance: (instance) => set({ selectedInstance: instance }),

  setLoading: (loading) => set({ loading }),

  setError: (error) => set({ error }),

  addInstance: (instance) => set((state) => ({ instances: [instance, ...state.instances] })),

  updateInstance: (id, updatedInstance) =>
    set((state) => ({
      instances: state.instances.map((inst) =>
        inst.id === id ? { ...inst, ...updatedInstance } : inst,
      ),
    })),

  removeInstance: (id) =>
    set((state) => ({
      instances: state.instances.filter((inst) => inst.id !== id),
    })),
}));