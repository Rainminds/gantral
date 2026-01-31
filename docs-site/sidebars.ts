import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';

/**
 * Sidebar ordering is intentional
 *
 * Reading flow:
 * Claim → Proof → Architecture → Governance → Adoption → Usage → Contribution
 */

const sidebars: SidebarsConfig = {
  docsSidebar: [
    // --- Root ---
    {
      type: 'doc',
      id: 'README',
      label: 'Welcome',
    },

    // --- Verifiability (PROOF FIRST) ---
    {
      type: 'category',
      label: 'Verifiability',
      link: {
        type: 'doc',
        id: 'verifiability/README',
      },
      items: [
        'verifiability/threat-model',
        'verifiability/commitment-artifact',
        'verifiability/replay-protocol',
        'verifiability/failure-semantics',
        'verifiability/security-model',
        'verifiability/non-claims',
      ],
    },

    // --- Positioning (Conceptual framing) ---
    {
      type: 'category',
      label: 'Positioning',
      items: [
        'positioning/ai-execution-control-plane-summary',
        'positioning/ai-execution-control-plane',
        'positioning/category-definition',
        'positioning/expansion-narrative',
        'positioning/what-is-gantral',
        'positioning/what-gantral-is-not',
      ],
    },

    // --- Architecture (Enforcement mechanics) ---
    {
      type: 'category',
      label: 'Architecture',
      items: [
        'architecture/trd',
        'architecture/authority-state-machine',
        'architecture/execution-authority-vs-agent-memory',
        'architecture/implementation-guide',
      ],
    },

    // --- Governance (Policy vs Authority) ---
    {
      type: 'category',
      label: 'Governance',
      items: [
        'governance/policy-vs-authority',
        'governance/oss-philosophy',
        'governance/contribution-model',
        'governance/license-commitment',
        'governance/cla',
        'governance/roadmap',
      ],
    },

    // --- Adoption (Operational rollout) ---
    {
      type: 'category',
      label: 'Adoption',
      items: [
        'adoption/README',
        'adoption/adoption-boundaries',
        'adoption/design-partners',
        'adoption/enterprise-onboarding',
      ],
    },

    // --- Guides (How to use & verify) ---
    {
      type: 'category',
      label: 'Guides',
      items: [
        'guides/auditor-verification',
        'guides/opa-integration',
        'guides/example-consumer-integration',
        'guides/demo',
      ],
    },

    // --- Product (Planning & scope) ---
    {
      type: 'category',
      label: 'Product',
      items: [
        'product/prd',
        'product/phase-wise-build-plan',
      ],
    },

    // --- Executive (Optional context) ---
    {
      type: 'category',
      label: 'Executive',
      items: ['executive/README'],
    },

    // --- Contributors ---
    {
      type: 'category',
      label: 'Contributing',
      items: [
        'contributors/getting-started',
        'contributors/how-to-contribute',
        'contributors/code-style-and-quality',
        'contributors/design-philosophy',
      ],
    },
  ],
};

export default sidebars;
