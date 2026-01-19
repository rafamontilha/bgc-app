/**
 * Lista de NCMs disponíveis para simulação
 * Baseado nos capítulos Onda 1: 02, 08, 84, 85
 *
 * Fonte: ComexStat - Principais produtos de exportação brasileira
 */

export interface NCMItem {
  code: string;
  description: string;
  chapter: string;
  chapterName: string;
}

/**
 * NCMs VALIDADOS para Simulação
 *
 * v0.4.0: Apenas 3 NCMs com dados reais no banco de dados
 *
 * IMPORTANTE: Esta lista contém APENAS NCMs que retornam dados reais do simulador.
 * Backend tem dados carregados para: 12010090, 17011400, 26011200
 */
export const AVAILABLE_NCMS: NCMItem[] = [
  // ========================================
  // Capítulo 12 - Sementes ✅ DISPONÍVEL (33 registros)
  // ========================================
  {
    code: '12010090',
    description: 'Outras sementes de soja, mesmo trituradas',
    chapter: '12',
    chapterName: 'Sementes e frutos oleaginosos',
  },

  // ========================================
  // Capítulo 17 - Açúcares ✅ DISPONÍVEL (33 registros)
  // ========================================
  {
    code: '17011400',
    description: 'Outros açúcares de cana, em bruto',
    chapter: '17',
    chapterName: 'Açúcares e produtos de confeitaria',
  },

  // ========================================
  // Capítulo 26 - Minérios ✅ DISPONÍVEL (20 registros)
  // ========================================
  {
    code: '26011200',
    description: 'Minérios de ferro aglomerados e seus concentrados',
    chapter: '26',
    chapterName: 'Minérios, escórias e cinzas',
  },
];

/**
 * Busca NCMs por código ou descrição
 */
export function searchNCMs(query: string): NCMItem[] {
  const normalizedQuery = query.toLowerCase().trim();

  if (!normalizedQuery) {
    return AVAILABLE_NCMS;
  }

  return AVAILABLE_NCMS.filter(
    (ncm) =>
      ncm.code.includes(normalizedQuery) ||
      ncm.description.toLowerCase().includes(normalizedQuery) ||
      ncm.chapterName.toLowerCase().includes(normalizedQuery)
  );
}

/**
 * Busca NCM por código exato
 */
export function getNCMByCode(code: string): NCMItem | undefined {
  return AVAILABLE_NCMS.find((ncm) => ncm.code === code);
}

/**
 * Formata NCM para exibição: "08051000 - Laranjas, frescas"
 */
export function formatNCM(ncm: NCMItem): string {
  return `${ncm.code} - ${ncm.description}`;
}
