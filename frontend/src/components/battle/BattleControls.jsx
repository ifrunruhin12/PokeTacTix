import PropTypes from 'prop-types';
import { motion } from 'framer-motion';

/**
 * BattleControls Component
 * Action buttons for battle: Attack, Defend, Pass, Sacrifice, Surrender
 * Disables buttons when not player's turn
 */
const BattleControls = ({ 
  onAttack, 
  onDefend, 
  onPass, 
  onSacrifice, 
  onSurrender,
  disabled = false,
  currentStamina = 0,
  maxStamina = 0,
  currentHp = 0,
  maxHp = 0,
  sacrificeCount = 0,
  playerSacrificeCost = null,
  className = ''
}) => {
  // Defend costs stamina (half of max HP + 1, matches backend GetDefendCost)
  const defendCost = Math.floor((maxHp + 1) / 2);
  const canDefend = currentStamina >= defendCost;

  // Next sacrifice HP cost comes from the backend (player_sacrifice_cost),
  // which is the single source of truth for the escalating 10/15/20 table.
  // It is null once the 3-sacrifice cap is reached.
  const sacrificeHpCost = playerSacrificeCost ?? null;
  const halfMaxStamina = Math.floor(maxStamina / 2);
  const canSacrifice = sacrificeHpCost !== null
    && currentHp > sacrificeHpCost
    && currentStamina < halfMaxStamina;

  const getSacrificeTooltip = () => {
    if (sacrificeCount >= 3) return 'Maximum sacrifices reached for this Pokémon (3 per battle)';
    if (sacrificeHpCost === null) return 'Sacrifice unavailable';
    if (currentHp <= sacrificeHpCost) return `Not enough HP (need more than ${sacrificeHpCost} HP)`;
    if (currentStamina >= halfMaxStamina) return `Stamina too high (must be below ${halfMaxStamina})`;
    return `Sacrifice ${sacrificeHpCost} HP to restore stamina (use ${sacrificeCount + 1}/3)`;
  };

  const buttons = [
    {
      label: 'Attack',
      icon: '⚔️',
      onClick: onAttack,
      color: 'from-red-600 to-red-700 hover:from-red-500 hover:to-red-600',
      disabled: disabled,
      tooltip: 'Choose a move to attack'
    },
    {
      label: 'Defend',
      icon: '🛡️',
      onClick: onDefend,
      color: 'from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600',
      disabled: disabled || !canDefend,
      tooltip: canDefend ? `Reduce damage (Cost: ${defendCost} stamina)` : `Not enough stamina (Need: ${defendCost})`
    },
    {
      label: 'Pass',
      icon: '⏭️',
      onClick: onPass,
      color: 'from-gray-600 to-gray-700 hover:from-gray-500 hover:to-gray-600',
      disabled: disabled,
      tooltip: 'Skip turn and recover stamina'
    },
    {
      label: 'Sacrifice',
      icon: '💥',
      onClick: onSacrifice,
      color: 'from-purple-600 to-purple-700 hover:from-purple-500 hover:to-purple-600',
      disabled: disabled || !canSacrifice,
      tooltip: getSacrificeTooltip()
    },
    {
      label: 'Surrender',
      icon: '🏳️',
      onClick: onSurrender,
      color: 'from-gray-700 to-gray-800 hover:from-gray-600 hover:to-gray-700',
      disabled: disabled,
      tooltip: 'Give up and end battle'
    }
  ];

  return (
    <div className={`flex flex-wrap gap-3 justify-center ${className}`}>
      {buttons.map((button, index) => (
        <motion.button
          key={index}
          whileHover={!button.disabled ? { scale: 1.05, y: -2 } : {}}
          whileTap={!button.disabled ? { scale: 0.95 } : {}}
          onClick={button.onClick}
          disabled={button.disabled}
          title={button.tooltip}
          className={`
            relative px-6 py-3 rounded-lg font-bold text-white shadow-lg
            transition-all duration-200
            ${button.disabled 
              ? 'bg-gray-700 opacity-50 cursor-not-allowed' 
              : `bg-gradient-to-br ${button.color} cursor-pointer`
            }
            min-w-[120px]
          `}
        >
          <div className="flex items-center justify-center gap-2">
            <span className="text-xl">{button.icon}</span>
            <span>{button.label}</span>
          </div>

          {/* Disabled overlay */}
          {button.disabled && (
            <div className="absolute inset-0 bg-black/20 rounded-lg" />
          )}
        </motion.button>
      ))}
    </div>
  );
};

BattleControls.propTypes = {
  onAttack: PropTypes.func.isRequired,
  onDefend: PropTypes.func.isRequired,
  onPass: PropTypes.func.isRequired,
  onSacrifice: PropTypes.func.isRequired,
  onSurrender: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
  currentStamina: PropTypes.number,
  maxStamina: PropTypes.number,
  currentHp: PropTypes.number,
  maxHp: PropTypes.number,
  sacrificeCount: PropTypes.number,
  playerSacrificeCost: PropTypes.number,
  className: PropTypes.string
};

export default BattleControls;
